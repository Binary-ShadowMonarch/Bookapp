<script lang="ts">
	// Imports for this page
	import { page } from '$app/stores';
	import NavBar from '$lib/NavBar.svelte';
	import { Book, Upload, Search, Bookmark, Play } from 'lucide-svelte';
</script>

<svelte:head>
	<title>Books | Your Digital Library for Reading and Organizing</title>
	<meta
		name="description"
		content="Discover, read, and organize your favorite books with our beautiful reading experience. Upload your own collection or explore thousands of free titles."
	/>
	<link rel="canonical" href={$page.url.href} />
</svelte:head>

<!-- Homepage -->
<NavBar>
	<div slot="name">Books</div>
	<div slot="right-buttons" class="flex justify-center gap-4 sm:flex-row">
		<a class="w-20" href="/login">
			<button
				aria-label="Login"
				class="flex cursor-pointer items-center gap-2 rounded-full px-3 py-1.5 text-sm font-semibold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 transition hover:bg-gray-100 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-gray-800"
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
					<path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" />
					<polyline points="10 17 15 12 10 7" />
					<line x1="15" x2="3" y1="12" y2="12" />
				</svg>
				<span class="sm:inline">Login</span>
			</button>
		</a>
		<a class="w-30" href="/signup">
			<button
				aria-label="Sign Up"
				class="flex cursor-pointer items-center gap-2 rounded-full bg-blue-500 bg-gradient-to-r from-blue-500 to-purple-600 px-3 py-1.5 text-sm font-semibold text-white hover:from-blue-600 hover:to-purple-700"
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
					<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
					<circle cx="9" cy="7" r="4" />
					<line x1="19" x2="19" y1="8" y2="14" />
					<line x1="22" x2="16" y1="11" y2="11" />
				</svg>
				<span class="sm:inline">Sign Up</span>
			</button>
		</a>
	</div>
</NavBar>

<main class="relative">
	<!-- Hero Section -->
	<section class="flex min-h-screen items-center justify-center px-4 pt-16 sm:px-6 lg:px-8">
		<div class="max-w- 3xl mx-auto text-center">
			<div class=" ">
				<h1
					class="mb-6 text-6xl font-bold leading-tight text-gray-900 sm:text-4xl lg:text-7xl dark:text-white"
				>
					Your Digital <span class="gradient-text">Library</span> Awaits
				</h1>
				<p class="mx-auto mb-8 max-w-2xl text-xl text-gray-700 sm:text-2xl dark:text-gray-300">
					Discover, read, and organize your favorite books with our beautiful reading experience
				</p>
				<div class="flex flex-col justify-center gap-4 sm:flex-row">
					<a
						href="/library"
						class="glow-effect hover-lift rounded-full bg-gradient-to-r from-blue-500 to-purple-600 px-8 py-4 text-lg font-semibold text-white transition-all duration-300 hover:from-blue-600 hover:to-purple-700"
					>
						Get Started Free
					</a>
					<button
						class="rounded-full border-2 border-gray-300 px-4 py-4 text-lg font-semibold text-gray-800 transition-all duration-300 hover:bg-gray-100 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-800"
					>
						<Play class="mr-2 inline h-5 w-5" />
						Watch Demo
					</button>
				</div>
			</div>

			<!-- Floating Book Animation -->
			<div class="relative mt-16">
				<div
					class="floating-card book-3d relative mx-auto h-40 w-32 overflow-hidden rounded-lg bg-gradient-to-br from-blue-400 to-purple-500 shadow-xl"
				>
					<div class="absolute inset-2 rounded-lg bg-white/10"></div>
					<div class="absolute bottom-4 left-4 right-4 h-1 rounded-full bg-white/20"></div>
					<div class="absolute bottom-7 left-4 right-4 h-1 rounded-full bg-white/20"></div>
					<div class="absolute bottom-10 left-4 right-4 h-1 rounded-full bg-white/20"></div>
				</div>
				<div class="pointer-events-none absolute inset-0 flex items-center justify-center">
					<div class="h-64 w-64 animate-pulse rounded-full bg-blue-500/10"></div>
				</div>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section
		class="
    px-4 py-20 sm:px-6 lg:px-8
  "
	>
		<div class="mx-auto max-w-7xl">
			<div class="mb-16 text-center">
				<h2 class="text- 3xl mb-6 font-bold text-gray-900 sm:text-5xl dark:text-white">
					Why Choose Books?
				</h2>
				<p class="mx-auto max-w-2xl text-xl text-gray-700 dark:text-gray-300">
					Everything you need for the perfect reading experience
				</p>
			</div>

			<div class="grid gap-8 md:grid-cols-3">
				<!-- Feature 1 -->
				<div class="acrylic hover-lift rounded-2xl p-8">
					<div
						class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-500 to-purple-600"
					>
						<Book class="h-8 w-8 text-white" />
					</div>
					<h3 class="mb-4 text-2xl font-bold text-gray-900 dark:text-white">
						Beautiful Reading UI
					</h3>
					<p class="leading-relaxed text-gray-700 dark:text-gray-300">
						Enjoy an immersive reading experience with our carefully crafted interface designed for
						comfort and focus.
					</p>
				</div>

				<!-- Feature 2 -->
				<div class="acrylic hover-lift rounded-2xl p-8">
					<div
						class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-green-500 to-teal-600"
					>
						<Upload class="h-8 w-8 text-white" />
					</div>
					<h3 class="mb-4 text-2xl font-bold text-gray-900 dark:text-white">Upload Your Books</h3>
					<p class="leading-relaxed text-gray-700 dark:text-gray-300">
						Import your own book collection and keep everything organized in one place. Support for
						multiple formats.
					</p>
				</div>

				<!-- Feature 3 -->
				<div class="acrylic hover-lift rounded-2xl p-8">
					<div
						class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-orange-500 to-red-600"
					>
						<Search class="h-8 w-8 text-white" />
					</div>
					<h3 class="mb-4 text-2xl font-bold text-gray-900 dark:text-white">Discover New Books</h3>
					<p class="leading-relaxed text-gray-700 dark:text-gray-300">
						Explore our curated collection of free books from authors worldwide. Find your next
						favorite read.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Free Library Section -->
	<section class="px-4 py-20 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-4xl">
			<div class="acrylic rounded-3xl p-8 text-center lg:p-12">
				<div class=" ">
					<div
						class="mb-8 inline-flex h-20 w-20 items-center justify-center rounded-full bg-gradient-to-br from-yellow-400 to-orange-500"
					>
						<Bookmark class="h-10 w-10 text-white" />
					</div>
					<h2 class="text- 3xl mb-6 font-bold text-gray-900 lg:text-5xl dark:text-white">
						Free Library Access
					</h2>
					<p class="mx-auto mb-8 max-w-2xl text-xl text-gray-700 dark:text-gray-300">
						Access thousands of free books from generous authors. Build your digital library without
						any cost.
					</p>
					<div class="flex flex-col justify-center gap-4 sm:flex-row">
						<div class="flex items-center space-x-3">
							<div class="pulse-glow h-2 w-2 rounded-full bg-green-400"></div>
							<span class="text-gray-700 dark:text-gray-300">Free to read</span>
						</div>
						<div class="flex items-center space-x-3">
							<div class="pulse-glow h-2 w-2 rounded-full bg-blue-400"></div>
							<span class="text-gray-700 dark:text-gray-300">No subscription required</span>
						</div>
						<div class="flex items-center space-x-3">
							<div class="pulse-glow h-2 w-2 rounded-full bg-purple-400"></div>
							<span class="text-gray-700 dark:text-gray-300">Offline reading</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section class="px-4 py-20 sm:px-6 lg:px-8">
		<div class="max-w- 3xl mx-auto text-center">
			<h2 class="text- 3xl mb-6 font-bold text-gray-900 lg:text-5xl dark:text-white">
				Ready to Start Reading?
			</h2>
			<p class="mb-8 text-xl text-gray-700 dark:text-gray-300">
				Join thousands of readers who've made Books their digital home
			</p>
			<a
				href="/library"
				class="glow-effect hover-lift inline-block rounded-full bg-gradient-to-r from-blue-500 to-purple-600 px-12 py-4 text-xl font-semibold text-white transition-all duration-300 hover:from-blue-600 hover:to-purple-700"
			>
				Get Started Free
			</a>
		</div>
	</section>
</main>

<!-- Footer -->
<footer
	class="border-t border-gray-200 bg-white/50 backdrop-blur-sm dark:border-gray-700 dark:bg-black/50"
>
	<div class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
		<div class="flex flex-col items-center justify-between md:flex-row">
			<div class="mb-4 flex items-center space-x-2 md:mb-0">
				<div
					class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-400 to-purple-500"
				>
					<Book class="h-5 w-5 text-white" />
				</div>
				<span class="text-xl font-bold text-gray-900 dark:text-white">Books</span>
			</div>
			<div class="flex space-x-6 text-gray-600 dark:text-gray-400">
				<a href="/" class="transition-colors hover:text-gray-900 dark:hover:text-white">Privacy</a>
				<a href="/" class="transition-colors hover:text-gray-900 dark:hover:text-white">Terms</a>
				<a href="/" class="transition-colors hover:text-gray-900 dark:hover:text-white">Support</a>
				<a href="/" class="transition-colors hover:text-gray-900 dark:hover:text-white">Contact</a>
			</div>
		</div>
		<div class="mt-8 border-t border-gray-200 pt-8 text-center dark:border-gray-700">
			<p class="text-gray-600 dark:text-gray-400">© 2024 Books. All rights reserved.</p>
		</div>
	</div>
</footer>

<style>
	.gradient-text {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	.acrylic {
		backdrop-filter: blur(10px);
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	:global(.dark) .acrylic {
		background: rgba(0, 0, 0, 0.2);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.floating-card {
		animation: float 6s ease-in-out infinite;
	}

	@keyframes float {
		0%,
		100% {
			transform: translateY(0px);
		}
		50% {
			transform: translateY(-10px);
		}
	}

	.glow-effect {
		box-shadow: 0 0 20px rgba(59, 130, 246, 0.3);
	}

	.hover-lift {
		transition: transform 0.3s ease;
	}

	.hover-lift:hover {
		transform: translateY(-5px);
	}

	.book-3d {
		transform-style: preserve-3d;
		transform: rotateX(15deg) rotateY(-15deg);
	}

	.pulse-glow {
		animation: pulse-glow 2s ease-in-out infinite;
	}
</style>
