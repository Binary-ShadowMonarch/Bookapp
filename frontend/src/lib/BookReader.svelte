<!-- src/lib/BookReader.svelte -->
<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { X, List, Settings, Loader2 } from 'lucide-svelte';
	import { fly } from 'svelte/transition';

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

	let darkMode = $state(true);
	let showChapterList = $state(false);
	let showSettings = $state(false);

	// Enhanced navigation state
	let isNavigating = $state(false);
	let isLoadingPrevious = $state(false);
	let isLoadingNext = $state(false);
	let navigationLock = $state(false);

	// --- DOM BINDINGS ---
	let readerContainer: HTMLElement;
	let settingsDropdown: HTMLElement;
	let loadPreviousSentinel: HTMLElement;
	let loadMoreSentinel: HTMLElement;

	// --- ASYNC & OBSERVERS ---
	let topObserver: IntersectionObserver;
	let bottomObserver: IntersectionObserver;
	let saveProgressTimeout: NodeJS.Timeout;
	let navigationTimeout: NodeJS.Timeout;

	// --- THEME DEFINITIONS ---
	const themes = {
		light: {
			body: { 'background-color': '#ffffff', color: '#111827' },
			a: { color: '#0000EE !important', 'text-decoration': 'underline' },
			'a:hover': { color: '#551A8B !important' }
		},
		dark: {
			body: { 'background-color': '#121212', color: '#e0e0e0' },
			a: { color: '#90cdf4 !important', 'text-decoration': 'underline' },
			'a:hover': { color: '#bee3f8 !important' }
		}
	};

	// --- LIFECYCLE & INITIALIZATION ---
	onMount(async () => {
		document.body.style.overflow = 'hidden';
		document.addEventListener('click', handleClickOutside);

		try {
			const ePub = (await import('epubjs')).default;
			const bookData = await loadBookData();
			const progressData = await loadProgress();

			book = ePub(bookData);

			rendition = book.renderTo(readerContainer, {
				manager: 'continuous',
				flow: 'scrolled-doc',
				width: '100%',
				height: '100%'
			});

			rendition.themes.register(themes);
			updateTheme();

			await rendition.display(progressData.location || undefined);
			await book.ready;
			chaptersForToc = book.navigation.toc;
			book.locations.generate(1600);

			setupEventHandlers();
			setupIntersectionObservers();

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
		clearTimeout(navigationTimeout);
		topObserver?.disconnect();
		bottomObserver?.disconnect();
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
			if (book.locations.length() > 0) {
				progress = book.locations.percentageFromCfi(currentLocation) * 100;
			}
			const currentChapter = book.navigation.get(location.start.href);
			if (currentChapter && currentChapter.label) {
				currentChapterLabel = currentChapter.label;
			}
			saveProgress();
		});

		rendition.on('rendered', () => {
			const view = rendition.manager.container;
			if (view) {
				view.prepend(loadPreviousSentinel);
				view.append(loadMoreSentinel);
			}
		});
	}

	function setupIntersectionObservers() {
		// Observer for loading PREVIOUS chapters
		topObserver = new IntersectionObserver(
			async ([entry]) => {
				if (
					!entry.isIntersecting ||
					navigationLock ||
					isLoadingPrevious ||
					rendition.location.start.index <= 0
				) {
					return;
				}

				try {
					navigationLock = true;
					isLoadingPrevious = true;

					// Add debounce to prevent rapid firing
					clearTimeout(navigationTimeout);
					navigationTimeout = setTimeout(async () => {
						try {
							const oldScrollHeight = readerContainer.scrollHeight;
							const oldScrollTop = readerContainer.scrollTop;

							await rendition.prev();

							// Maintain scroll position
							requestAnimationFrame(() => {
								const newScrollHeight = readerContainer.scrollHeight;
								const heightDiff = newScrollHeight - oldScrollHeight;
								if (heightDiff > 0) {
									readerContainer.scrollTop = oldScrollTop + heightDiff;
								}
							});
						} catch (e) {
							console.error('Error loading previous section:', e);
						} finally {
							isLoadingPrevious = false;
							// Release lock after a short delay to prevent immediate re-triggering
							setTimeout(() => {
								navigationLock = false;
							}, 300);
						}
					}, 100);
				} catch (e) {
					console.error('Error in top observer:', e);
					isLoadingPrevious = false;
					navigationLock = false;
				}
			},
			{
				root: readerContainer,
				rootMargin: '50px 0px 0px 0px' // Trigger slightly before reaching the edge
			}
		);

		// Observer for loading NEXT chapters
		bottomObserver = new IntersectionObserver(
			async ([entry]) => {
				if (!entry.isIntersecting || navigationLock || isLoadingNext) {
					return;
				}

				try {
					navigationLock = true;
					isLoadingNext = true;

					// Add debounce
					clearTimeout(navigationTimeout);
					navigationTimeout = setTimeout(async () => {
						try {
							await rendition.next();
						} catch (e) {
							console.error('Error loading next section:', e);
						} finally {
							isLoadingNext = false;
							// Release lock after a short delay
							setTimeout(() => {
								navigationLock = false;
							}, 300);
						}
					}, 100);
				} catch (e) {
					console.error('Error in bottom observer:', e);
					isLoadingNext = false;
					navigationLock = false;
				}
			},
			{
				root: readerContainer,
				rootMargin: '0px 0px 50px 0px' // Trigger slightly before reaching the edge
			}
		);

		topObserver.observe(loadPreviousSentinel);
		bottomObserver.observe(loadMoreSentinel);
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
		}, 1500);
	}

	function goToChapter(href: string) {
		navigationLock = true;
		rendition.display(href);
		showChapterList = false;
		setTimeout(() => {
			navigationLock = false;
		}, 500);
	}

	function toggleDarkMode() {
		darkMode = !darkMode;
		updateTheme();
	}

	function updateTheme() {
		const themeName = darkMode ? 'dark' : 'light';
		rendition?.themes.select(themeName);
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
		class="flex flex-shrink-0 items-center justify-between border-b border-gray-200 bg-white px-4 py-3 dark:border-gray-700 dark:bg-gray-800"
	>
		<div class="flex min-w-0 items-center gap-2">
			<button
				onclick={onClose}
				class="flex-shrink-0 rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
				aria-label="Close reader"
			>
				<X class="h-5 w-5 text-gray-600 dark:text-gray-300" />
			</button>
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
			>
				<List class="h-5 w-5 text-gray-600 dark:text-gray-300" />
			</button>
			<div class="relative" bind:this={settingsDropdown}>
				<button
					onclick={() => (showSettings = !showSettings)}
					class="rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
					aria-label="Settings"
				>
					<Settings class="h-5 w-5 text-gray-600 dark:text-gray-300" />
				</button>
				{#if showSettings}
					<div
						transition:fly={{ y: -5, duration: 150 }}
						class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-lg border bg-white shadow-lg dark:border-gray-700 dark:bg-gray-800"
					>
						<div class="flex items-center justify-between px-4 py-3">
							<span class="text-sm font-medium text-gray-700 dark:text-gray-300">
								{darkMode ? 'Dark Mode' : 'Light Mode'}
							</span>
							<button
								aria-label="Dark mode toggle"
								onclick={toggleDarkMode}
								class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 {darkMode
									? 'bg-blue-600'
									: 'bg-gray-200'}"
								role="switch"
								aria-checked={darkMode}
							>
								<span
									class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform {darkMode
										? 'translate-x-6'
										: 'translate-x-1'}"
								></span>
							</button>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</header>

	<div class="h-1 w-full bg-gray-200 dark:bg-gray-700">
		<div class="h-full bg-blue-500 transition-all duration-500" style:width="{progress}%"></div>
	</div>

	<div class="relative flex flex-1 overflow-hidden">
		{#if showChapterList && chaptersForToc.length > 0}
			<aside
				transition:fly={{ x: -20, duration: 200 }}
				class="w-72 flex-shrink-0 border-r bg-white dark:border-gray-700 dark:bg-gray-800"
			>
				<h3 class="p-4 text-lg font-semibold dark:text-white">Table of Contents</h3>
				<div class="h-[calc(100%-4rem)] overflow-y-auto">
					{#each chaptersForToc as chapter}
						<button
							onclick={() => goToChapter(chapter.href)}
							class="block w-full px-4 py-2 text-left text-sm transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
						>
							{chapter.label.trim()}
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
					</div>
				</div>
			{/if}

			<!-- Loading indicators for infinite scroll -->
			{#if isLoadingPrevious}
				<div class="absolute left-1/2 top-0 z-10 -translate-x-1/2 transform">
					<div
						class="flex items-center gap-2 rounded-full bg-black/80 px-3 py-2 text-white backdrop-blur-sm"
					>
						<Loader2 class="h-4 w-4 animate-spin" />
						<span class="text-sm">Loading previous...</span>
					</div>
				</div>
			{/if}

			{#if isLoadingNext}
				<div class="absolute bottom-0 left-1/2 z-10 -translate-x-1/2 transform">
					<div
						class="flex items-center gap-2 rounded-full bg-black/80 px-3 py-2 text-white backdrop-blur-sm"
					>
						<Loader2 class="h-4 w-4 animate-spin" />
						<span class="text-sm">Loading next...</span>
					</div>
				</div>
			{/if}

			<div bind:this={loadPreviousSentinel} class="h-px w-full"></div>

			<div
				bind:this={readerContainer}
				class="reader-view h-full w-full transition-opacity duration-300"
				style:opacity={isLoading ? 0 : 1}
			></div>

			<div bind:this={loadMoreSentinel} class="h-px w-full"></div>
		</main>
	</div>
</div>

<style>
	.reader-view {
		overflow-y: scroll;
		-webkit-overflow-scrolling: touch;
		scrollbar-width: thin;
		scrollbar-color: #a0a0a0 #f0f0f0;
	}

	:global(.epub-view > iframe) {
		background-color: transparent !important;
	}
	:global(.epub-arrow) {
		display: none;
	}
</style>
