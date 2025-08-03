<script lang="ts">
	import BookCard from '../../lib/BookCard.svelte';
	import NavBar from '../../lib/NavBar.svelte';
	import SearchBar from '../../lib/SearchBar.svelte';
	import EpubUpload from '../../lib/EpubUpload.svelte';
	import BookReader from '../../lib/BookReader.svelte';
	import { epubParser } from '$lib/epub-parser-pool';
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { LogOut, Settings, HelpCircle, User, Store } from 'lucide-svelte';

	let { data }: { data: { user: { email: string } } } = $props();

	// --- TYPE DEFINITIONS ---
	type BookStatus = 'read' | 'unread' | 'finished';
	const modes = ['all', 'read', 'unread', 'finished'] as const;
	type Mode = (typeof modes)[number];

	interface Book {
		id: string; // Ensure this is a stable, unique ID from your backend
		title: string;
		author: string;
		image: string;
		status: BookStatus;
		completion: number;
		fileUrl: string;
	}

	// --- STATE ---
	let books: Book[] = $state([]);
	let searchTerm = $state('');
	let statusFilter: Mode = $state('all');
	let showMenu = $state(false);
	let showReader = $state(false);
	let currentBookId = $state('');
	let currentBookUrl = $state('');
	let isLoading = $state(false); // ✅ Start as false to prevent flash
	let isAuthenticating = $state(true);
	let isAuthenticated = $state(false);
	let booksLoaded = $state(false);

	// --- LIFECYCLE & DATA FETCHING ---
	onMount(async () => {
		const isAuth = await checkAuth();
		if (isAuth && !booksLoaded) {
			// Only load if not already loaded
			await loadExistingBooks();
			booksLoaded = true;
		}
	});

	async function checkAuth(): Promise<boolean> {
		try {
			let res = await fetch('/api/protected/profile', { method: 'GET', credentials: 'include' });
			if (res.status === 401) {
				const refresh = await fetch('/api/refresh', { method: 'GET', credentials: 'include' });
				if (refresh.ok) {
					res = await fetch('/api/protected/profile', { method: 'GET', credentials: 'include' });
				}
			}
			if (!res.ok) {
				isAuthenticated = false;
				goto('/login');
				return false;
			}
			isAuthenticated = true;
			return true;
		} catch {
			isAuthenticated = false;
			goto('/login');
			return false;
		} finally {
			// ✅ This ensures we stop showing the primary loader after the check.
			isAuthenticating = false;
		}
	}

	async function parseEpubMetadata(file: File) {
		try {
			const result = await epubParser.parseEpub(file);

			// If coverUrl is a blob URL, convert it to data URL
			if (result.coverUrl && result.coverUrl.startsWith('blob:')) {
				try {
					const response = await fetch(result.coverUrl);
					const blob = await response.blob();
					const dataUrl = await new Promise<string>((resolve) => {
						const reader = new FileReader();
						reader.onload = () => resolve(reader.result as string);
						reader.readAsDataURL(blob);
					});
					result.coverUrl = dataUrl;
				} catch (error) {
					console.error('Failed to convert blob URL to data URL:', error);
					result.coverUrl = null;
				}
			}

			return result;
		} catch (error) {
			console.error('Failed to parse EPUB:', error);
			return {
				title: file.name.replace('.epub', '') || 'Unknown Title',
				author: 'Unknown Author',
				coverUrl: null
			};
		}
	}

	async function loadExistingBooks() {
		// Only load if not already loaded
		if (booksLoaded) return;

		isLoading = true;
		try {
			const res = await fetch('/api/protected/library', { credentials: 'include' });
			if (!res.ok) throw new Error('Failed to fetch library');

			const { files } = await res.json();
			if (!files || files.length === 0) {
				books = [];
				return;
			}

			// Use a Map to track existing books and prevent duplicates
			const existingBooks = new Map(books.map((book) => [book.id, book]));
			const newBooks: Book[] = [];

			for (const fileInfo of files) {
				try {
					// Skip if book already exists
					if (existingBooks.has(fileInfo.id)) continue;

					const [blobRes, progressData] = await Promise.all([
						fetch(fileInfo.url, { credentials: 'include' }),
						loadBookProgress(fileInfo.id)
					]);

					if (!blobRes.ok) continue;

					const blob = await blobRes.blob();
					const file = new File([blob], fileInfo.name, { type: blob.type });
					const meta = await parseEpubMetadata(file);

					newBooks.push({
						id: fileInfo.id,
						title: meta.title,
						author: meta.author,
						image: meta.coverUrl || '/candle.webp',
						status:
							progressData.completion === 100
								? 'finished'
								: progressData.completion > 0
									? 'read'
									: 'unread',
						completion: progressData.completion || 0,
						fileUrl: fileInfo.url
					});
				} catch (e) {
					console.error(`Error processing book ${fileInfo.name}:`, e);
				}
			}

			books = [...books, ...newBooks]; // Add only new books
		} catch (e) {
			console.error('Error loading existing books:', e);
		} finally {
			isLoading = false;
		}
	}
	async function loadBookProgress(bookId: string) {
		try {
			const response = await fetch(`/api/protected/library/progress?bookId=${bookId}`, {
				credentials: 'include'
			});
			if (response.ok) {
				const data = await response.json();
				return {
					completion: data.progress || 0,
					location: data.location || ''
				};
			}
		} catch (error) {
			console.error('Error loading book progress:', error);
		}
		return { completion: 0, location: '' };
	}

	// --- EVENT HANDLERS ---
	function handleBookAdded(book: Book) {
		// When a book is added, also make sure loading is false.
		isLoading = false;
		books = [...books, book];
	}

	function handleOpenBook(bookId: string, fileUrl: string) {
		currentBookId = bookId;
		currentBookUrl = fileUrl;
		showReader = true;
	}

	function handleCloseReader() {
		// Don't reset immediately - wait for animation
		setTimeout(() => {
			showReader = false;
			currentBookId = '';
			currentBookUrl = '';
		}, 300); // Match the transition duration
	}

	function handleProgressUpdate(bookId: string, progress: number) {
		books = books.map((book) => {
			if (book.id === bookId) {
				return {
					...book,
					completion: progress,
					status: progress === 100 ? 'finished' : progress > 0 ? 'read' : 'unread'
				};
			}
			return book;
		});
	}

	async function logout() {
		await fetch('/api/logout', { method: 'POST', credentials: 'include' });
		goto('/login');
	}

	// --- DERIVED STATE ---
	let filteredBooks = $derived(
		books.filter((b) => {
			const q = searchTerm.trim().toLowerCase();
			const matchesSearch =
				!q || b.title.toLowerCase().includes(q) || b.author.toLowerCase().includes(q);
			const matchesStatus = statusFilter === 'all' || b.status === statusFilter;
			return matchesSearch && matchesStatus;
		})
	);
</script>

<svelte:head>
	<title>Books | Library</title>
	<meta
		name="description"
		content="Discover, read, and organize your favorite books with our beautiful reading experience. Upload your own collection or explore thousands of free titles."
	/>
</svelte:head>

{#if isAuthenticating}
	<div class="flex min-h-screen flex-col items-center justify-center text-white">
		<div class="loading loading-ring loading-lg"></div>
		<p class="mt-4 text-lg">Authenticating...</p>
	</div>
{:else if !isAuthenticated}
	<div class="flex min-h-screen flex-col items-center justify-center text-white">
		<div class="loading loading-ring loading-lg"></div>
		<p class="mt-4 text-lg">Authenticating......</p>
	</div>
{:else}
	{#if showReader}
		<BookReader
			bookId={currentBookId}
			fileUrl={currentBookUrl}
			onClose={handleCloseReader}
			onProgressUpdate={handleProgressUpdate}
		/>
	{/if}

	<div hidden={showReader}>
		<NavBar showSecondaryRow={true}>
			<div class="flex justify-center" slot="name">Library</div>
			<div slot="center" class="relative flex items-center">
				<SearchBar bind:value={searchTerm} />
			</div>
			<div slot="right-buttons" class="flex justify-center gap-4 sm:flex-row">
				<button
					aria-label="Store"
					class="ml-2 flex cursor-pointer items-center gap-2 rounded-full bg-indigo-500 px-3 py-1.5 text-sm font-semibold text-white shadow-sm transition hover:bg-indigo-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
					><Store /><span class="hidden sm:inline">Store</span></button
				><EpubUpload onBookAdded={handleBookAdded} />
				<div class="relative">
					<button
						onclick={() => (showMenu = !showMenu)}
						class="rounded-full p-2 transition hover:bg-gray-200 dark:hover:bg-gray-700"
						aria-label="User menu"><User class="h-6 w-6 text-gray-800 dark:text-white" /></button
					>
					{#if showMenu}
						<div
							class="absolute right-0 z-50 mt-2 w-48 rounded-lg bg-white text-sm shadow-lg dark:bg-gray-800"
						>
							<button
								class="flex w-full items-center gap-2 px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700"
								><Settings class="h-4 w-4" /> Settings</button
							><button
								class="flex w-full items-center gap-2 px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700"
								><HelpCircle class="h-4 w-4" /> Help & Support</button
							>
							<hr class="border-gray-300 dark:border-gray-600" />
							<button
								onclick={logout}
								class="flex w-full items-center gap-2 px-4 py-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900"
								><LogOut class="h-4 w-4" /> Logout</button
							>
						</div>
					{/if}
				</div>
			</div>
			<div slot="secondary-center-buttons" class="flex gap-2">
				{#each modes as mode}
					<button
						type="button"
						onclick={() => (statusFilter = mode)}
						class="          cursor-pointer rounded-full bg-black/20 px-4 py-1 text-sm font-medium text-white backdrop-blur-sm transition-colors duration-200 ease-in-out hover:bg-white/15 {statusFilter ===
						mode
							? 'bg-white/30'
							: ''}        "
						>{mode === 'all' ? 'All Books' : mode[0].toUpperCase() + mode.slice(1)}</button
					>
				{/each}
			</div>
		</NavBar>
		<div class="flex min-h-screen flex-col items-center justify-center py-4 pt-32">
			<div class="mx-auto w-full px-2 md:px-4">
				{#if isLoading}
					<div class="flex flex-col items-center justify-center space-y-4 text-center">
						<div class="loading loading-ring loading-lg text-white"></div>
						<p class="text-lg text-gray-300">Loading your library...</p>
					</div>
				{:else if filteredBooks.length > 0}
					<div
						class="grid grid-cols-2 gap-x-4 gap-y-6 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
					>
						{#each filteredBooks as book (book.id)}
							<BookCard {...book} onOpenBook={handleOpenBook} />
						{/each}
					</div>
				{:else}
					<div class="flex h-full flex-col items-center justify-center pt-24 text-center">
						<p class="text-xl text-gray-400">Upload or hit the store to get started.</p>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
