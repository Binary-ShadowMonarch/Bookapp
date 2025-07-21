<!-- src/lib/BookReader.svelte -->
<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { X, List, Settings } from 'lucide-svelte';

	interface Props {
		bookId: string;
		fileUrl?: string;
		arrayBuffer?: ArrayBuffer;
		onClose: () => void;
		onProgressUpdate: (bookId: string, progress: number, location: string) => void;
	}

	let { bookId, fileUrl, arrayBuffer, onClose, onProgressUpdate }: Props = $props();

	let readerContainer: HTMLElement;
	let book: any;
	let rendition: any;
	let isLoading = $state(true);
	let currentLocation = $state('');
	let progress = $state(0);
	let showChapterList = $state(false);
	let darkMode = $state(false);
	let chapters: any[] = $state([]);
	let currentChapter = $state('');
	let showSettings = $state(false);

	// Concurrent loading states
	let bookReady = $state(false);
	let locationsGenerated = $state(false);
	let progressLoaded = $state(false);

	onMount(async () => {
		try {
			document.body.style.overflow = 'hidden';

			// Start all async operations concurrently
			const [bookData, progressData] = await Promise.all([loadBookData(), loadProgress()]);

			// Initialize book and rendition
			const ePub = (await import('epubjs')).default;
			book = ePub(bookData);

			rendition = book.renderTo(readerContainer, {
				manager: 'continuous',
				flow: 'scrolled-doc',
				width: '100%',
				height: '100%'
			});

			updateTheme();

			// Display immediately (fastest priority)
			if (currentLocation) {
				await rendition.display(currentLocation);
			} else {
				await rendition.display();
			}

			// Book is now displayable
			isLoading = false;

			// Start background operations concurrently
			Promise.all([initializeBookMetadata(), generateLocations(), setupEventHandlers()]).catch(
				console.error
			);
		} catch (error) {
			console.error('Error loading book:', error);
			isLoading = false;
		}
	});

	// Separate function for book data loading
	async function loadBookData(): Promise<ArrayBuffer> {
		if (arrayBuffer) {
			return arrayBuffer;
		}

		if (fileUrl) {
			const response = await fetch(fileUrl, { credentials: 'include' });
			if (!response.ok) {
				throw new Error(`Failed to fetch book: ${response.statusText}`);
			}
			return response.arrayBuffer();
		}

		throw new Error('No book data or URL provided');
	}

	// Enhanced progress loading with caching
	async function loadProgress(): Promise<void> {
		try {
			const response = await fetch(
				`http://localhost:8080/protected/library/progress?bookId=${bookId}`,
				{ credentials: 'include' }
			);
			if (response.ok) {
				const data = await response.json();
				progress = data.progress || 0;
				currentLocation = data.location || '';
			}
		} catch (error) {
			console.log('Progress API not available - using defaults:', error);
		} finally {
			progressLoaded = true;
		}
	}

	// Background metadata initialization
	async function initializeBookMetadata(): Promise<void> {
		try {
			await book.ready;
			chapters = book.navigation.toc.map((item: any) => ({
				id: item.id,
				href: item.href,
				label: item.label.trim()
			}));
			bookReady = true;
		} catch (error) {
			console.error('Error loading book metadata:', error);
		}
	}

	// Background location generation
	async function generateLocations(): Promise<void> {
		try {
			await book.locations.generate(1600);
			locationsGenerated = true;
			// Recalculate progress now that locations are available
			if (currentLocation && book.locations.length() > 0) {
				progress = Math.round(book.locations.percentageFromCfi(currentLocation) * 100);
			}
		} catch (error) {
			console.error('Error generating locations:', error);
		}
	}

	// Setup event handlers
	async function setupEventHandlers(): Promise<void> {
		if (rendition) {
			rendition.on('relocated', handleLocationChange);
		}
	}

	// Debounced save progress for better performance
	let saveProgressTimeout: NodeJS.Timeout;
	async function saveProgress(): Promise<void> {
		if (!bookId) return;

		// Debounce saves to avoid excessive API calls
		clearTimeout(saveProgressTimeout);
		saveProgressTimeout = setTimeout(async () => {
			try {
				await fetch('http://localhost:8080/protected/library/progress', {
					method: 'POST',
					credentials: 'include',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ bookId, progress, location: currentLocation })
				});
				onProgressUpdate(bookId, progress, currentLocation);
			} catch (error) {
				console.log('Progress save API not available:', error);
				onProgressUpdate(bookId, progress, currentLocation);
			}
		}, 500); // 500ms debounce
	}

	function handleLocationChange(location: any): void {
		currentLocation = location.start.cfi;

		// Calculate progress if locations are ready
		if (locationsGenerated && book.locations?.length() > 0) {
			progress = Math.round(book.locations.percentageFromCfi(location.start.cfi) * 100);
		}

		// Update current chapter
		const spine = book.spine.get(location.start.cfi);
		if (spine && chapters.length > 0) {
			const chapter = chapters.find((ch) => ch.href === spine.href);
			if (chapter) {
				currentChapter = chapter.label;
			}
		}

		saveProgress();
	}

	function goToChapter(href: string): void {
		rendition?.display(href);
		showChapterList = false;
	}

	function toggleDarkMode(): void {
		darkMode = !darkMode;
		updateTheme();
	}

	function updateTheme(): void {
		if (!rendition) return;
		const theme = darkMode
			? { body: { 'background-color': '#1f2937', color: '#f9fafb' } }
			: { body: { 'background-color': '#ffffff', color: '#111827' } };
		rendition.themes.default(theme);
	}

	function handleGlobalKeydown(event: KeyboardEvent): void {
		if (event.key === 'Escape') {
			onClose();
		}
	}

	onDestroy(() => {
		document.body.style.overflow = 'auto';
		if (saveProgressTimeout) {
			clearTimeout(saveProgressTimeout);
		}
		if (rendition) {
			rendition.destroy();
		}
	});
</script>

<svelte:window on:keydown={handleGlobalKeydown} />

<div class="fixed inset-0 z-50 flex flex-col bg-white dark:bg-gray-900">
	<header
		class="flex items-center justify-between border-b border-gray-200 bg-white px-4 py-3 dark:border-gray-700 dark:bg-gray-800"
	>
		<div class="flex items-center gap-2">
			<button
				onclick={onClose}
				class="rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
				aria-label="Close reader"
			>
				<X class="h-5 w-5 text-gray-600 dark:text-gray-300" />
			</button>
			<div class="truncate text-sm text-gray-700 dark:text-gray-200" title={currentChapter}>
				{currentChapter || 'Reading...'}
			</div>
		</div>

		<div class="flex items-center gap-2">
			<div class="text-sm font-medium text-gray-600 dark:text-gray-300">{progress}%</div>
			<button
				onclick={() => (showChapterList = !showChapterList)}
				class="rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700 {chapters.length === 0
					? 'cursor-not-allowed opacity-50'
					: ''}"
				aria-label="Table of contents"
				disabled={chapters.length === 0}
			>
				<List class="h-5 w-5 text-gray-600 dark:text-gray-300" />
			</button>
			<div class="relative">
				<button
					onclick={() => (showSettings = !showSettings)}
					class="rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
					aria-label="Settings"
				>
					<Settings class="h-5 w-5 text-gray-600 dark:text-gray-300" />
				</button>
				{#if showSettings}
					<div
						class="absolute right-0 z-10 mt-2 w-48 rounded-lg border bg-white shadow-lg dark:border-gray-700 dark:bg-gray-800"
					>
						<button
							onclick={() => {
								toggleDarkMode();
								showSettings = false;
							}}
							class="flex w-full items-center gap-3 px-3 py-2 text-left text-sm hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
						>
							{darkMode ? '☀️ Light Mode' : '🌙 Dark Mode'}
						</button>
					</div>
				{/if}
			</div>
		</div>
	</header>

	<div class="h-1 w-full bg-gray-200 dark:bg-gray-700">
		<div class="h-full bg-blue-500 transition-all" style="width: {progress}%;"></div>
	</div>

	<div class="relative flex flex-1 overflow-hidden">
		{#if showChapterList && chapters.length > 0}
			<aside class="w-72 flex-shrink-0 border-r bg-white dark:border-gray-700 dark:bg-gray-800">
				<h3 class="p-4 text-lg font-semibold dark:text-white">Table of Contents</h3>
				<div class="h-[calc(100%-4rem)] overflow-y-auto">
					{#each chapters as chapter}
						<button
							onclick={() => goToChapter(chapter.href)}
							class="block w-full px-4 py-2 text-left text-sm hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
						>
							{chapter.label}
						</button>
					{/each}
				</div>
			</aside>
		{/if}

		<main class="relative flex-1">
			{#if isLoading}
				<div class="flex h-full items-center justify-center">
					<div class="flex flex-col items-center gap-4">
						<div
							class="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-500"
						></div>
						<p class="text-gray-600 dark:text-gray-300">Loading book...</p>
						{#if progressLoaded}
							<p class="text-xs text-gray-500">Progress loaded ✓</p>
						{/if}
						{#if bookReady}
							<p class="text-xs text-gray-500">Metadata loaded ✓</p>
						{/if}
						{#if locationsGenerated}
							<p class="text-xs text-gray-500">Locations generated ✓</p>
						{/if}
					</div>
				</div>
			{/if}
			<div bind:this={readerContainer} class="h-full w-full"></div>
		</main>
	</div>
</div>
