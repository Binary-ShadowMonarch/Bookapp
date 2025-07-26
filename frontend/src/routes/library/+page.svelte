<script lang="ts">
	import BookCard from '../../lib/BookCard.svelte';
	import NavBar from '../../lib/NavBar.svelte';
	import SearchBar from '../../lib/SearchBar.svelte';
	import EpubUpload from '../../lib/EpubUpload.svelte';
	import BookReader from '../../lib/BookReader.svelte';
	import { epubParser } from '$lib/epub-parser-pool';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { LogOut, Settings, HelpCircle, User2 } from 'lucide-svelte';

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
	let isLoading = $state(true);

	// --- LIFECYCLE & DATA FETCHING ---
	onMount(async () => {
		const ok = await checkAuth();
		if (ok) {
			await loadExistingBooks();
		} else {
			// Ensure loading stops if auth fails
			isLoading = false;
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
				goto('/login');
				return false;
			}
			return true;
		} catch {
			goto('/login');
			return false;
		}
	}

	async function parseEpubMetadata(file: File) {
		try {
			return await epubParser.parseEpub(file);
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
		isLoading = true;
		try {
			const res = await fetch('/api/protected/library', { credentials: 'include' });
			if (!res.ok) throw new Error('Failed to fetch library');

			const { files } = await res.json();

			// ✅ FIX: The problematic early return has been removed.
			// The logic now gracefully handles an empty `files` array,
			// allowing the `finally` block to always run.

			const bookPromises = (files || []).map(
				async (fileInfo: { id: string; name: string; url: string }) => {
					try {
						const [blobRes, progressData] = await Promise.all([
							fetch(fileInfo.url, { credentials: 'include' }),
							loadBookProgress(fileInfo.id)
						]);

						if (!blobRes.ok) return null;

						const blob = await blobRes.blob();
						const file = new File([blob], fileInfo.name, { type: blob.type });
						const meta = await parseEpubMetadata(file);

						return {
							id: fileInfo.id,
							title: meta.title,
							author: meta.author,
							image: meta.coverUrl || '/default.webp',
							status:
								progressData.completion === 100
									? 'finished'
									: progressData.completion > 0
										? 'read'
										: 'unread',
							completion: progressData.completion || 0,
							fileUrl: fileInfo.url
						};
					} catch (e) {
						console.error(`Error processing book ${fileInfo.name}:`, e);
						return null;
					}
				}
			);

			const results = await Promise.all(bookPromises);
			books = results.filter((book): book is Book => book !== null);
		} catch (e) {
			console.error('Error loading existing books:', e);
			books = [];
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
		showReader = false;
		currentBookId = '';
		currentBookUrl = '';
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

{#if showReader}
	<BookReader
		bookId={currentBookId}
		fileUrl={currentBookUrl}
		onClose={handleCloseReader}
		onProgressUpdate={handleProgressUpdate}
	/>
{/if}

<NavBar showSecondaryRow={true}
	><div class="flex justify-center" slot="name">Library</div>
	<div slot="center" class="relative flex items-center"><SearchBar bind:value={searchTerm} /></div>
	<div slot="right-buttons" class="flex justify-center gap-4 sm:flex-row">
		<button
			aria-label="Store"
			class="flex cursor-pointer items-center gap-2 rounded-full bg-indigo-500 px-3 py-1.5 text-sm font-semibold text-white shadow-sm transition hover:bg-indigo-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
			><svg
				xmlns="http://www.w3.org/2000/svg"
				width="20"
				height="20"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
				class="h-5 w-5"
				><path d="M6 2 3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" /><path d="M3 6h18" /><path
					d="M16 10a4 4 0 0 1-8 0"
				/></svg
			><span class="hidden sm:inline">Store</span></button
		><EpubUpload onBookAdded={handleBookAdded} />
		<div class="relative">
			<button
				onclick={() => (showMenu = !showMenu)}
				class="rounded-full p-2 transition hover:bg-gray-200 dark:hover:bg-gray-700"
				aria-label="User menu"><User2 class="h-6 w-6 text-gray-800 dark:text-white" /></button
			>{#if showMenu}<div
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
				</div>{/if}
		</div>
	</div>
	<div slot="secondary-center-buttons" class="flex gap-2">
		{#each modes as mode}<button
				type="button"
				onclick={() => (statusFilter = mode)}
				class="          cursor-pointer rounded-full bg-black/20 px-4 py-1 text-sm font-medium text-white backdrop-blur-sm transition-colors duration-200 ease-in-out hover:bg-white/15 {statusFilter ===
				mode
					? 'bg-white/30'
					: ''}        "
				>{mode === 'all' ? 'All Books' : mode[0].toUpperCase() + mode.slice(1)}</button
			>{/each}
	</div></NavBar
>
<div class="flex min-h-screen flex-col items-center justify-center py-4 pt-32">
	<div class="mx-auto w-full px-2 md:px-4">
		{#if isLoading}<div class="flex flex-col items-center justify-center space-y-4 text-center">
				<div class="loading loading-ring loading-lg text-white"></div>
				<p class="text-lg text-gray-300">Loading your library...</p>
			</div>{:else if filteredBooks.length > 0}<div
				class="grid grid-cols-2 gap-x-4 gap-y-6 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#each filteredBooks as book (book.id)}<BookCard
						{...book}
						onOpenBook={handleOpenBook}
					/>{/each}
			</div>{:else}<div class="flex h-full flex-col items-center justify-center pt-24 text-center">
				<p class="text-xl text-gray-400">Upload or hit the store to get started.</p>
			</div>{/if}
	</div>
</div>
