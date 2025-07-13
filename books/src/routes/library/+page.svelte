<script lang="ts">
	import BookCard from '../../lib/BookCard.svelte';
	import NavBar from '../../lib/NavBar.svelte';
	import SearchBar from '../../lib/SearchBar.svelte';
	import EpubUpload from '../../lib/EpubUpload.svelte';
	import { onMount } from 'svelte';

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
	}

	// Initialize with mock books
	let books: Book[] = [
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'To Kill a Mockingbird',
		// 	author: 'Harper Lee',
		// 	image: '/default.webp',
		// 	status: 'read',
		// 	completion: 100
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: '1984',
		// 	author: 'George Orwell',
		// 	image: '/default.webp',
		// 	status: 'finished',
		// 	completion: 100
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'Pride and Prejudice',
		// 	author: 'Jane Austen',
		// 	image: '/default.webp',
		// 	status: 'unread',
		// 	completion: 0
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'The Great Gatsby',
		// 	author: 'F. Scott Fitzgerald',
		// 	image: '/default.webp',
		// 	status: 'read',
		// 	completion: 80
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'Moby-Dick',
		// 	author: 'Herman Melville',
		// 	image: '/default.webp',
		// 	status: 'unread',
		// 	completion: 10
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'War and Peace',
		// 	author: 'Leo Tolstoy',
		// 	image: '/default.webp',
		// 	status: 'finished',
		// 	completion: 100
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'The Catcher in the Rye',
		// 	author: 'J.D. Salinger',
		// 	image: '/default.webp',
		// 	status: 'read',
		// 	completion: 60
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'The Hobbit',
		// 	author: 'J.R.R. Tolkien',
		// 	image: '/default.webp',
		// 	status: 'finished',
		// 	completion: 100
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'Brave New World',
		// 	author: 'Aldous Huxley',
		// 	image: '/default.webp',
		// 	status: 'unread',
		// 	completion: 0
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'The Odyssey',
		// 	author: 'Homer',
		// 	image: '/default.webp',
		// 	status: 'read',
		// 	completion: 50
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: '1984',
		// 	author: 'George Orwell',
		// 	image: '/default.webp',
		// 	status: 'finished',
		// 	completion: 100
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'The Hobbit',
		// 	author: 'J.R.R. Tolkien',
		// 	image: '/default.webp',
		// 	status: 'finished',
		// 	completion: 100
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'To Kill a Mockingbird',
		// 	author: 'Harper Lee',
		// 	image: '/default.webp',
		// 	status: 'read',
		// 	completion: 100
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'Moby-Dick',
		// 	author: 'Herman Melville',
		// 	image: '/default.webp',
		// 	status: 'unread',
		// 	completion: 10
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'Pride and Prejudice',
		// 	author: 'Jane Austen',
		// 	image: '/default.webp',
		// 	status: 'unread',
		// 	completion: 0
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'The Odyssey',
		// 	author: 'Homer',
		// 	image: '/default.webp',
		// 	status: 'read',
		// 	completion: 50
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'War and Peace',
		// 	author: 'Leo Tolstoy',
		// 	image: '/default.webp',
		// 	status: 'finished',
		// 	completion: 100
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'The Catcher in the Rye',
		// 	author: 'J.D. Salinger',
		// 	image: '/default.webp',
		// 	status: 'read',
		// 	completion: 60
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'Brave New World',
		// 	author: 'Aldous Huxley',
		// 	image: '/default.webp',
		// 	status: 'unread',
		// 	completion: 0
		// },
		// {
		// 	id: crypto.randomUUID(),
		// 	title: 'The Great Gatsby',
		// 	author: 'F. Scott Fitzgerald',
		// 	image: '/default.webp',
		// 	status: 'read',
		// 	completion: 80
		// }
	];

	let searchTerm = '';
	let statusFilter: Mode = 'all';

	// Load books from localStorage on mount


	function handleBookAdded(book: Book) {
		books = [...books, book];
	}

	$: filteredBooks = books.filter((b) => {
		const q = searchTerm.trim().toLowerCase();
		const matchesSearch =
			!q || b.title.toLowerCase().includes(q) || b.author.toLowerCase().includes(q);

		const matchesStatus = statusFilter === 'all' || b.status === statusFilter;

		return matchesSearch && matchesStatus;
	});
</script>

<!-- Nav bar -->
<NavBar showSecondaryRow={true}>
	<div class="flex justify-center" slot="name">Library</div>

	<!-- search -->
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

		<!-- Replace the upload button with the EpubUpload component -->
		<EpubUpload onBookAdded={handleBookAdded} />
	</div>

	<!-- secondary row buttons -->
	<div slot="secondary-center-buttons" class="flex gap-2">
		{#each modes as mode}
			<button
				type="button"
				on:click={() => (statusFilter = mode)}
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

<!-- Library page -->
<div class="flex min-h-screen flex-col items-center justify-center py-4 pt-20">
	<div class="mx-auto w-full px-2 md:px-4">
		<div
			class="grid grid-cols-2 gap-x-4 gap-y-6 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
		>
			{#if filteredBooks.length}
				{#each filteredBooks as book (book.id)}
					<BookCard {...book} />
				{/each}
			{:else}
				<p class="col-span-full text-center text-gray-500">
					Add a new book "{searchTerm}" {statusFilter !== 'all' ? `and status=${statusFilter}` : ''}
				</p>
			{/if}
		</div>
	</div>
</div>
