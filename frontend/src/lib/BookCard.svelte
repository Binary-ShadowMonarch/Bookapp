<!-- src/lib/BookCard.svelte -->
<script lang="ts">
	export let id: string;
	export let title: string;
	export let author: string;
	export let image: string = '/default.webp';
	export let status: 'read' | 'unread' | 'finished' = 'unread';
	export let completion: number = 0;
	export let fileUrl: string;
	export let onOpenBook: (bookId: string, fileUrl: string) => void;

	function handleImgError(event: Event) {
		const target = event.target as HTMLImageElement | null;
		if (target) target.src = '/default.webp';
	}

	function openBook() {
		onOpenBook(id, fileUrl);
	}
</script>

<div
	class="group relative flex aspect-[2/3] min-w-[140px] max-w-[180px] cursor-pointer flex-col justify-end overflow-hidden rounded-2xl bg-white/60 shadow-lg backdrop-blur-md transition-transform hover:scale-105 sm:min-w-[160px] sm:max-w-[220px] md:min-w-[180px] md:max-w-[240px] lg:min-w-[200px] lg:max-w-[260px] dark:bg-black/40"
>
	<img
		src={image}
		alt={title}
		class="absolute inset-0 z-0 h-full w-full object-cover object-center transition group-hover:brightness-90"
		on:error={handleImgError}
	/>

	<!-- Open button overlay - shows on hover -->
	<div
		class="absolute inset-0 z-10 flex items-center justify-center bg-black/50 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
	>
		<button
			on:click={openBook}
			class="rounded-full bg-white/90 px-6 py-3 font-semibold text-gray-900 shadow-lg transition-transform hover:scale-105 hover:bg-white"
		>
			Open Book
		</button>
	</div>

	<!-- Info background -->
	<div
		class="relative z-20 mt-auto flex w-full flex-col rounded-b-2xl bg-white/90 p-3 pb-4 shadow-md dark:bg-black/80"
	>
		<h3 class="truncate text-base font-semibold text-gray-900 dark:text-gray-100">{title}</h3>
		<p class="truncate text-sm text-gray-700 dark:text-gray-300">{author}</p>
		<span
			class="mt-2 inline-flex items-center gap-2 rounded-full border px-2 py-0.5 text-xs font-medium shadow-sm {status ===
			'read'
				? 'border-green-300 bg-green-200 text-green-800 dark:border-green-600 dark:bg-green-700 dark:text-green-100'
				: ''} {status === 'unread'
				? 'border-yellow-300 bg-yellow-200 text-yellow-800 dark:border-yellow-600 dark:bg-yellow-700 dark:text-yellow-100'
				: ''} {status === 'finished'
				? 'border-blue-300 bg-blue-200 text-blue-800 dark:border-blue-600 dark:bg-blue-700 dark:text-blue-100'
				: ''}"
		>
			{status.charAt(0).toUpperCase() + status.slice(1)}
			<span class="block h-1 w-12 overflow-hidden rounded-full bg-gray-300 dark:bg-gray-700">
				<span
					class="block h-1 rounded-full transition-all duration-300"
					class:bg-green-500={completion === 100}
					class:bg-blue-500={completion < 100 && completion > 0}
					class:bg-gray-400={completion === 0}
					style="width: {completion}%"
				></span>
			</span>
			<span class="ml-1 font-semibold">{completion}%</span>
		</span>
	</div>
</div>
