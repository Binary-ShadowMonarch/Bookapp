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
	type BookStatus = 'read' | 'unread' | 'finished';
	const modes = ['all', 'read', 'unread', 'finished'] as const;
	type Mode = (typeof modes)[number];

	interface Book {
		id: string;
		title: string;
		author: string;
		image: string;
		status: BookStatus;
		completion: number;
		fileUrl: string;
	}

	let books: Book[] = $state([]);
	let searchTerm = $state('');
	let statusFilter: Mode = $state('all');
	let showMenu = $state(false);
	let showReader = $state(false);
	let currentBookId = $state('');
	let currentBookUrl = $state('');
	let loading = $state(true);
	// 2) tweak your check() so it never reloads the page (we’ll control it)
	async function check(): Promise<boolean> {
		let res = await fetch('http://localhost:8080/protected/library', {
			method: 'GET',
			credentials: 'include'
		});

		if (res.status === 401) {
			const refresh = await fetch('http://localhost:8080/refresh', {
				method: 'GET',
				credentials: 'include'
			});
			if (refresh.ok) {
				// retry once
				res = await fetch('http://localhost:8080/protected/library', {
					method: 'GET',
					credentials: 'include'
				});
			}
		}

		if (res.status === 401 || !res.ok) {
			goto('/login');
			return false;
		}
		return true;
	}

	onMount(async () => {
		// 3) start spinner
		loading = true;

		const ok = await check();
		if (!ok) {
			return;
		}
		await loadExistingBooks();
		loading = false;
	});

	// Replace existing parseEpubMetadata function
	async function parseEpubMetadata(file: File) {
		try {
			return await epubParser.parseEpub(file);
		} catch (error) {
			console.error('Failed to parse EPUB:', error);
			// Fallback metadata
			return {
				id: crypto.randomUUID(),
				title: file.name.replace('.epub', '') || 'Unknown Title',
				author: 'Unknown Author',
				coverUrl: null
			};
		}
	}
	// Optimized loadExistingBooks with batch processing
	async function loadExistingBooks() {
		try {
			const res = await fetch('http://localhost:8080/protected/library', {
				credentials: 'include'
			});
			if (!res.ok) return;

			const { files } = await res.json();

			// Fetch all blobs concurrently
			const filePromises = files.map(async (fileInfo: { name: string; url: string }) => {
				const blobRes = await fetch(fileInfo.url, { credentials: 'include' });
				if (!blobRes.ok) return null;

				const blob = await blobRes.blob();
				return {
					file: new File([blob], fileInfo.name, { type: blob.type }),
					url: fileInfo.url
				};
			});
			const fileResults = await Promise.allSettled(filePromises);
			const validFiles = fileResults
				.filter(
					(result): result is PromiseFulfilledResult<{ file: File; url: string } | null> =>
						result.status === 'fulfilled' && result.value !== null
				)
				.map((result) => result.value!);

			// Parse all files concurrently using worker pool
			const metadataPromises = validFiles.map(({ file, url }) =>
				parseEpubMetadata(file).then((meta) => ({ ...meta, fileUrl: url }))
			);

			const progressPromises = validFiles.map(
				({ file }) =>
					// Generate consistent ID for progress lookup
					loadBookProgress(crypto.randomUUID()) // You'll need to fix this ID generation
			);

			const [metadataResults, progressResults] = await Promise.all([
				Promise.allSettled(metadataPromises),
				Promise.allSettled(progressPromises)
			]);

			const newBooks: Book[] = [];
			metadataResults.forEach((result, index) => {
				if (result.status === 'fulfilled') {
					const meta = result.value;
					const progressResult = progressResults[index];
					const progressData =
						progressResult.status === 'fulfilled'
							? progressResult.value
							: { completion: 0, status: 'unread' };

					const book: Book = {
						id: meta.id,
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
						fileUrl: meta.fileUrl
					};
					newBooks.push(book);
				}
			});

			books = newBooks;
		} catch (e) {
			console.error('Error loading existing EPUBs:', e);
		}
	}

	async function loadBookProgress(bookId: string) {
		try {
			const response = await fetch(
				`http://localhost:8080/protected/library/progress?bookId=${bookId}`,
				{ credentials: 'include' }
			);
			if (response.ok) {
				const data = await response.json();
				return {
					completion: data.progress || 0,
					status: data.progress === 100 ? 'finished' : data.progress > 0 ? 'read' : 'unread'
				};
			}
		} catch (error) {
			console.error('Error loading book progress:', error);
		}
		return { completion: 0, status: 'unread' as BookStatus };
	}

	function handleBookAdded(book: Book) {
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

	let filteredBooks = $derived(
		books.filter((b) => {
			const q = searchTerm.trim().toLowerCase();
			const matchesSearch =
				!q || b.title.toLowerCase().includes(q) || b.author.toLowerCase().includes(q);
			const matchesStatus = statusFilter === 'all' || b.status === statusFilter;
			return matchesSearch && matchesStatus;
		})
	);

	function toggleMenu() {
		showMenu = !showMenu;
	}

	async function logout() {
		await fetch('http://localhost:8080/logout', {
			method: 'POST',
			credentials: 'include'
		});
		await goto('/login');
	}
</script>

{#if loading}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30">
		<span class="loading loading-ring loading-xl"></span>
	</div>
{/if}

{#if showReader}
	<BookReader
		bookId={currentBookId}
		fileUrl={currentBookUrl}
		onClose={handleCloseReader}
		onProgressUpdate={handleProgressUpdate}
	/>
{/if}
{#if !loading}
	<NavBar showSecondaryRow={true}>
		<div class="flex justify-center" slot="name">Library</div>

		<div slot="center" class="relative flex items-center">
			<SearchBar bind:value={searchTerm} />
		</div>

		<div slot="right-buttons" class="flex justify-center gap-4 sm:flex-row">
			<button
				aria-label="Store"
				class="flex cursor-pointer items-center gap-2 rounded-full bg-indigo-500 px-3 py-1.5 text-sm font-semibold text-white shadow-sm transition hover:bg-indigo-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
			>
				<svg
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
				>
					<path d="M6 2 3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
					<path d="M3 6h18" />
					<path d="M16 10a4 4 0 0 1-8 0" />
				</svg>
				<span class="hidden sm:inline">Store</span>
			</button>

			<EpubUpload onBookAdded={handleBookAdded} />

			<div class="relative">
				<button
					onclick={toggleMenu}
					class="rounded-full p-2 transition hover:bg-gray-200 dark:hover:bg-gray-700"
					aria-label="User menu"
				>
					<User2 class="h-6 w-6 text-gray-800 dark:text-white" />
				</button>

				{#if showMenu}
					<div
						class="absolute right-0 z-50 mt-2 w-48 rounded-lg bg-white text-sm shadow-lg dark:bg-gray-800"
					>
						<button
							class="flex w-full items-center gap-2 px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700"
						>
							<Settings class="h-4 w-4" /> Settings
						</button>
						<button
							class="flex w-full items-center gap-2 px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700"
						>
							<HelpCircle class="h-4 w-4" /> Help & Support
						</button>
						<hr class="border-gray-300 dark:border-gray-600" />
						<button
							onclick={logout}
							class="flex w-full items-center gap-2 px-4 py-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900"
						>
							<LogOut class="h-4 w-4" /> Logout
						</button>
					</div>
				{/if}
			</div>
		</div>

		<div slot="secondary-center-buttons" class="flex gap-2">
			{#each modes as mode}
				<button
					type="button"
					onclick={() => (statusFilter = mode)}
					class="
          cursor-pointer rounded-full bg-black/20 px-4 py-1 text-sm
          font-medium text-white
          backdrop-blur-sm transition-colors duration-200
          ease-in-out
          hover:bg-white/15
          {statusFilter === mode ? 'bg-white/30' : ''}
        "
				>
					{mode === 'all' ? 'All Books' : mode[0].toUpperCase() + mode.slice(1)}
				</button>
			{/each}
		</div>
	</NavBar>

	<div class="flex min-h-screen flex-col items-center justify-center py-4 pt-20">
		<div class="mx-auto w-full px-2 md:px-4">
			<div
				class="grid grid-cols-2 gap-x-4 gap-y-6 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#if filteredBooks.length}
					{#each filteredBooks as book (book.id)}
						<BookCard {...book} onOpenBook={handleOpenBook} />
					{/each}
				{:else}
					<p class="col-span-full text-center text-gray-500">
						Add a new book "{searchTerm}" {statusFilter !== 'all'
							? `and status=${statusFilter}`
							: ''}
					</p>
				{/if}
			</div>
		</div>
	</div>
{/if}
