<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { X, List, Settings } from 'lucide-svelte';

	// --- PROPS & INTERFACES ---
	interface Props {
		bookId: string;
		fileUrl?: string;
		onClose: () => void;
		onProgressUpdate: (bookId: string, progress: number, location: string) => void;
	}
	let { bookId, fileUrl, onClose, onProgressUpdate }: Props = $props();

	// --- STATE ---
	let book: any;
	let rendition: any;
	let chaptersForToc: any[] = $state([]);

	let isLoading = $state(true);
	let currentLocation = $state('');
	let currentChapterLabel = $state('Reading...');
	let progress = $state(0);

	let darkMode = $state(true); // Default to dark mode
	let showChapterList = $state(false);
	let showSettings = $state(false);

	let readerContainer: HTMLElement; // The single container for epub.js
	let loadMoreSentinel: HTMLElement; // The trigger to load the next chapter
	let settingsDropdown: HTMLElement; // Reference to settings dropdown

	let intersectionObserver: IntersectionObserver;
	let saveProgressTimeout: NodeJS.Timeout;

	// --- LIFECYCLE & INITIALIZATION ---
	onMount(async () => {
		document.body.style.overflow = 'hidden';
		
		// Add click outside listener for settings
		document.addEventListener('click', handleClickOutside);
		
		try {
			const ePub = (await import('epubjs')).default;
			const bookData = await loadBookData();
			const progressData = await loadProgress();

			book = ePub(bookData);

			// Render the book into our single container
			rendition = book.renderTo(readerContainer, {
				manager: 'continuous',
				flow: 'scrolled-doc',
				width: '100%',
				height: '100%'
			});

			// Set the theme before displaying
			updateTheme();

			// Display the book, jumping to the last saved location if it exists
			await rendition.display(progressData.location || undefined);

			await book.ready;
			chaptersForToc = book.navigation.toc;

			// Generate locations in the background for accurate progress
			book.locations.generate(1600);

			setupEventHandlers();
			setupIntersectionObserver();

			isLoading = false;
		} catch (error) {
			console.error('Error loading book:', error);
			isLoading = false;
		}
	});

	onDestroy(() => {
		document.body.style.overflow = 'auto';
		document.removeEventListener('click', handleClickOutside);
		clearTimeout(saveProgressTimeout);
		intersectionObserver?.disconnect();
		book?.destroy();
	});

	// --- CORE LOGIC & EVENT HANDLERS ---
	async function loadBookData(): Promise<ArrayBuffer> {
		if (!fileUrl) throw new Error('No book URL provided');
		const response = await fetch(fileUrl, { credentials: 'include' });
		if (!response.ok) throw new Error(`Failed to fetch book: ${response.statusText}`);
		return response.arrayBuffer();
	}

	async function loadProgress() {
		try {
			const res = await fetch(`/api/protected/library/progress?bookId=${bookId}`, {
				credentials: 'include'
			});
			if (res.ok) return await res.json();
		} catch (error) {
			console.error('Could not load progress:', error);
		}
		return { progress: 0, location: '' };
	}

	function setupEventHandlers() {
		rendition.on('relocated', (location: any) => {
			currentLocation = location.start.cfi;

			// Update progress percentage using epub.js's location generation
			if (book.locations.length() > 0) {
				progress = book.locations.percentageFromCfi(currentLocation) * 100;
			}

			// Update the chapter title in the header
			const currentChapter = book.navigation.get(location.start.href);
			if (currentChapter && currentChapter.label) {
				currentChapterLabel = currentChapter.label;
			}

			saveProgress();
		});

		// When a new chapter is added, this event fires
		rendition.on('rendered', (section: any) => {
			// Move the sentinel to the end of the new content
			section.output.append(loadMoreSentinel);
		});
	}

	function setupIntersectionObserver() {
		intersectionObserver = new IntersectionObserver(
			([entry]) => {
				// When the sentinel is on screen, load the next chapter
				if (entry.isIntersecting) {
					rendition.next();
				}
			},
			{ root: readerContainer }
		);

		intersectionObserver.observe(loadMoreSentinel);
	}

	function saveProgress() {
		clearTimeout(saveProgressTimeout);
		saveProgressTimeout = setTimeout(() => {
			fetch('/api/protected/library/progress', {
				method: 'POST',
				credentials: 'include',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ bookId, progress: Math.round(progress), location: currentLocation })
			});
			onProgressUpdate(bookId, Math.round(progress), currentLocation);
		}, 1000);
	}

	function goToChapter(href: string) {
		rendition.display(href);
		showChapterList = false;
	}

	function toggleDarkMode() {
		darkMode = !darkMode;
		updateTheme();
	}

	function updateTheme() {
		const theme = darkMode
			? { 'background-color': '#1f2937', color: '#f9fafb' }
			: { 'background-color': '#ffffff', color: '#111827' };
		rendition?.themes.override('color', theme['color']);
		rendition?.themes.override('background-color', theme['background-color']);
	}

	function handleClickOutside(event: MouseEvent) {
		if (showSettings && settingsDropdown && !settingsDropdown.contains(event.target as Node)) {
			showSettings = false;
		}
	}
</script>

<svelte:window on:keydown={(e) => e.key === 'Escape' && onClose()} />

<div class="fixed inset-0 z-50 flex flex-col bg-white dark:bg-gray-900">
	<header
		class="flex items-center justify-between border-b border-gray-200 bg-white px-4 py-3 dark:border-gray-700 dark:bg-gray-800"
	>
		<div class="flex min-w-0 items-center gap-2">
			<button
				onclick={onClose}
				class="flex-shrink-0 rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
				aria-label="Close reader"><X class="h-5 w-5 text-gray-600 dark:text-gray-300" /></button
			>
			<div class="truncate text-sm text-gray-700 dark:text-gray-200" title={currentChapterLabel}>
				{currentChapterLabel}
			</div>
		</div>
		<div class="flex flex-shrink-0 items-center gap-2">
			<div class="text-sm font-medium text-gray-600 dark:text-gray-300">
				{Math.round(progress)}%
			</div>
			<button
				onclick={() => (showChapterList = !showChapterList)}
				class="rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
				aria-label="Table of contents"
				><List class="h-5 w-5 text-gray-600 dark:text-gray-300" /></button
			>
			<div class="relative" bind:this={settingsDropdown}>
				<button
					onclick={() => (showSettings = !showSettings)}
					class="rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
					aria-label="Settings"
					><Settings class="h-5 w-5 text-gray-600 dark:text-gray-300" /></button
				>
				{#if showSettings}
					<div
						class="absolute right-0 z-10 mt-2 w-56 rounded-lg border bg-white shadow-lg dark:border-gray-700 dark:bg-gray-800"
					>
						<div class="flex items-center justify-between px-4 py-3">
							<span class="text-sm font-medium text-gray-700 dark:text-gray-300">
								{darkMode ? 'Dark Mode' : 'Light Mode'}
							</span>
							<button
								onclick={toggleDarkMode}
								class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 {darkMode ? 'bg-blue-600' : 'bg-gray-200'}"
								role="switch"
								aria-checked={darkMode}
							>
								<span
									class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform {darkMode ? 'translate-x-6' : 'translate-x-1'}"
								></span>
							</button>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</header>

	<div class="h-1 w-full bg-gray-200 dark:bg-gray-700">
		<div class="h-full bg-blue-500 transition-all" style:width="{progress}%"></div>
	</div>

	<div class="relative flex flex-1 overflow-hidden">
		{#if showChapterList && chaptersForToc.length > 0}
			<aside class="w-72 flex-shrink-0 border-r bg-white dark:border-gray-700 dark:bg-gray-800">
				<h3 class="p-4 text-lg font-semibold dark:text-white">Table of Contents</h3>
				<div class="h-[calc(100%-4rem)] overflow-y-auto">
					{#each chaptersForToc as chapter}
						<button
							onclick={() => goToChapter(chapter.href)}
							class="block w-full px-4 py-2 text-left text-sm hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
							>{chapter.label}</button
						>
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
					</div>
				</div>
			{/if}

			<div
				bind:this={readerContainer}
				class="reader-view h-full w-full"
				style:opacity={isLoading ? 0 : 1}
			></div>

			<div bind:this={loadMoreSentinel} class="h-px"></div>
		</main>
	</div>
</div>

<style>
	/* Ensures the epub.js container allows scrolling */
	.reader-view {
		overflow-y: scroll;
		-webkit-overflow-scrolling: touch; /* Smooth scrolling on iOS */
	}
	/* Hides the default epub.js navigation arrows */
	:global(.epub-view > iframe) {
		background-color: transparent !important;
	}
	:global(.epub-arrow) {
		display: none !important;
	}
</style>