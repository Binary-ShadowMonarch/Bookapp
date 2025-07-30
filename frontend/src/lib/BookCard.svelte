<!-- this is a book card component that shows a single book in the library -->
<!-- it displays the cover image, title, author, and reading progress -->
<script lang="ts">
	// these are the props that get passed in from the parent component
	export let id: string;
	export let title: string;
	export let author: string;
	export let image: string = '/default.webp'; // default cover if none provided
	export let status: 'read' | 'unread' | 'finished' = 'unread';
	export let completion: number = 0; // reading progress percentage
	export let fileUrl: string; // URL to the actual book file
	export let onOpenBook: (bookId: string, fileUrl: string) => void; // callback when book is opened

	// if the cover image fails to load, show the default image instead
	function handleImgError(event: Event) {
		console.log('DEBUG: Book cover image failed to load, using default');
		const target = event.target as HTMLImageElement | null;
		if (target) target.src = '/default.webp';
	}

	// when someone clicks to open the book, call the parent's function
	function openBook() {
		console.log('DEBUG: Opening book:', id, title);
		onOpenBook(id, fileUrl);
	}
</script>

<!-- the main book card container -->
<!-- it has a nice hover effect and works in both light and dark mode -->
<div
	class="group relative flex aspect-[2/3] min-w-[140px] max-w-[180px] cursor-pointer flex-col justify-end overflow-hidden rounded-2xl bg-white/60 shadow-lg backdrop-blur-md transition-transform hover:scale-105 sm:min-w-[160px] sm:max-w-[220px] md:min-w-[180px] md:max-w-[240px] lg:min-w-[200px] lg:max-w-[260px] dark:bg-black/40"
>
	<!-- the book cover image -->
	<!-- it fills the entire card and gets slightly darker on hover -->
	<img
		src={image}
		alt={title}
		class="absolute inset-0 z-0 h-full w-full object-cover object-center transition group-hover:brightness-90"
		on:error={handleImgError}
	/>

	<!-- this overlay appears when you hover over the book -->
	<!-- it shows an "Open Book" button -->
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

	<!-- the info section at the bottom of the card -->
	<!-- shows title, author, status, and reading progress -->
	<div
		class="relative z-20 mt-auto flex w-full flex-col rounded-b-2xl bg-white/90 p-3 pb-4 shadow-md dark:bg-black/80"
	>
		<!-- book title - gets truncated if too long -->
		<h3 class="truncate text-base font-semibold text-gray-900 dark:text-gray-100">{title}</h3>
		<!-- author name - also gets truncated -->
		<p class="truncate text-sm text-gray-700 dark:text-gray-300">{author}</p>
		
		<!-- status badge with reading progress bar -->
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
			<!-- progress bar background -->
			<span class="block h-1 w-12 overflow-hidden rounded-full bg-gray-300 dark:bg-gray-700">
				<!-- the actual progress bar that fills based on completion percentage -->
				<span
					class="block h-1 rounded-full transition-all duration-300"
					class:bg-green-500={completion === 100}
					class:bg-blue-500={completion < 100 && completion > 0}
					class:bg-gray-400={completion === 0}
					style="width: {completion}%"
				></span>
			</span>
			<!-- completion percentage text -->
			<span class="ml-1 font-semibold">{completion}%</span>
		</span>
	</div>
</div>
