<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { X, List, Settings } from 'lucide-svelte';

	interface Props {
		bookId: string;
		fileUrl: string;
		onClose: () => void;
		onProgressUpdate: (bookId: string, progress: number, location: string) => void;
	}

	let { bookId, fileUrl, onClose, onProgressUpdate }: Props = $props();

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

	onMount(async () => {
		try {
			const response = await fetch(fileUrl, { credentials: 'include' });
			if (!response.ok) {
				throw new Error(`Failed to fetch book: ${response.statusText}`);
			}
			const arrayBuffer = await response.arrayBuffer();
			document.body.style.overflow = 'hidden';

			const ePub = (await import('epubjs')).default;
			book = ePub(arrayBuffer);
			await book.ready;

			chapters = book.navigation.toc.map((item: any) => ({
				id: item.id,
				href: item.href,
				label: item.label.trim()
			}));

			await loadProgress();

			// --- KEY CHANGE HERE ---
			// Use the 'continuous' manager for seamless scrolling
			rendition = book.renderTo(readerContainer, {
				manager: 'continuous', // This enables true infinite scroll
				flow: 'scrolled-doc',
				width: '100%',
				height: '100%'
			});

			updateTheme();

			await book.locations.generate(1600);

			if (currentLocation) {
				await rendition.display(currentLocation);
			} else {
				await rendition.display();
			}

			rendition.on('relocated', handleLocationChange);
			isLoading = false;
		} catch (error) {
			console.error('Error loading book:', error);
			isLoading = false;
		}
	});

	onDestroy(() => {
		document.body.style.overflow = 'auto';
		if (rendition) {
			rendition.destroy();
		}
	});

	// The manual scroll handler is no longer needed.

	async function loadProgress() {
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
		}
	}

	async function saveProgress() {
		if (!bookId) return;
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
	}

	function handleLocationChange(location: any) {
		currentLocation = location.start.cfi;
		if (book.locations && book.locations.length() > 0) {
			progress = Math.round(book.locations.percentageFromCfi(location.start.cfi) * 100);
		}
		const spine = book.spine.get(location.start.cfi);
		if (spine) {
			const chapter = chapters.find((ch) => ch.href === spine.href);
			if (chapter) {
				currentChapter = chapter.label;
			}
		}
		saveProgress();
	}

	function goToChapter(href: string) {
		rendition?.display(href);
		showChapterList = false;
	}

	function toggleDarkMode() {
		darkMode = !darkMode;
		updateTheme();
	}

	function updateTheme() {
		if (!rendition) return;
		const theme = darkMode
			? { body: { 'background-color': '#1f2937', color: '#f9fafb' } }
			: { body: { 'background-color': '#ffffff', color: '#111827' } };
		rendition.themes.default(theme);
	}

	function handleGlobalKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			onClose();
		}
	}
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
				class="rounded-lg p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
				aria-label="Table of contents"
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
		{#if showChapterList}
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
					<p class="text-gray-600 dark:text-gray-300">Loading book...</p>
				</div>
			{/if}
			<div bind:this={readerContainer} class="h-full w-full"></div>
		</main>
	</div>
</div>
