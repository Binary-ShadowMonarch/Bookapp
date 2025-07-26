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
    <div slot = "name">Books</div>
	<div slot = "right-buttons" class="flex sm:flex-row gap-4 justify-center">
        <a href = "/login">
            <button aria-label="Login"
            class="flex items-center gap-2 rounded-full px-3 py-1.5 text-sm font-semibold text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 transition hover:bg-gray-100 dark:text-gray-300 dark:ring-gray-600 dark:hover:bg-gray-800">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5">
                    <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" />
                    <polyline points="10 17 15 12 10 7" />
                    <line x1="15" x2="3" y1="12" y2="12" />
                </svg>
                <span class="hidden sm:inline">Login</span>
            </button>
        </a>
        <a href = "/signup">

        <button aria-label="Sign Up" class="flex items-center gap-2 rounded-full bg-blue-600 px-3 py-1.5 text-sm font-semibold text-white shadow-sm transition hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5" >
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                <circle cx="9" cy="7" r="4" />
                <line x1="19" x2="19" y1="8" y2="14" />
                <line x1="22" x2="16" y1="11" y2="11" />
            </svg>
            <span class="hidden sm:inline">Sign Up</span>
        </button>
    </a>
	</div>
</NavBar>

<main class="relative">
	<!-- Hero Section -->
	<section class="min-h-screen flex items-center justify-center px-4 sm:px-6 lg:px-8 pt-16  ">
		<div class="max-w- 3xl mx-auto text-center">
			<div class=" ">
				<h1 class="text-6xl sm:text-4xl lg:text-7xl font-bold mb-6 leading-tight text-gray-900 dark:text-white">
					Your Digital <span class="gradient-text">Library</span> Awaits
				</h1>
				<p class="text-xl sm:text-2xl text-gray-700 dark:text-gray-300 mb-8 max-w-2xl mx-auto">
					Discover, read, and organize your favorite books with our beautiful reading experience
				</p>
				<div class="flex flex-col sm:flex-row gap-4 justify-center">
					<a href="/library" class="px-8 py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-full font-semibold text-lg hover:from-blue-600 hover:to-purple-700 transition-all duration-300 glow-effect hover-lift">
						Get Started Free
					</a>
					<button class="px-4 py-4 border-2 border-gray-300 dark:border-gray-600 rounded-full font-semibold text-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-300 text-gray-800 dark:text-gray-200">
						<Play class="w-5 h-5 inline mr-2" />
						Watch Demo
					</button>
				</div>
			</div>
			
			<!-- Floating Book Animation -->
			<div class="mt-16 relative">
				<div class="floating-card book-3d mx-auto w-32 h-40 bg-gradient-to-br from-blue-400 to-purple-500 rounded-lg shadow-xl relative overflow-hidden">
					<div class="absolute inset-2 bg-white/10 rounded-lg"></div>
					<div class="absolute bottom-4 left-4 right-4 h-1 bg-white/20 rounded-full"></div>
					<div class="absolute bottom-7 left-4 right-4 h-1 bg-white/20 rounded-full"></div>
					<div class="absolute bottom-10 left-4 right-4 h-1 bg-white/20 rounded-full"></div>
				</div>
				<div class="absolute inset-0 flex items-center justify-center pointer-events-none">
					<div class="w-64 h-64 bg-blue-500/10 rounded-full animate-pulse"></div>
				</div>
			</div>
		</div>
	</section>
	
	<!-- Features Section -->
<section
  class="
    py-20 px-4 sm:px-6 lg:px-8
  ">
		<div class="max-w-7xl mx-auto">
			<div class="text-center mb-16  ">
				<h2 class="text- 3xl sm:text-5xl font-bold mb-6 text-gray-900 dark:text-white">Why Choose Books?</h2>
				<p class="text-xl text-gray-700 dark:text-gray-300 max-w-2xl mx-auto">
					Everything you need for the perfect reading experience
				</p>
			</div>
			
			<div class="grid md:grid-cols-3 gap-8">
				<!-- Feature 1 -->
				<div class="acrylic rounded-2xl p-8 hover-lift  ">
					<div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mb-6">
						<Book class="w-8 h-8 text-white" />
					</div>
					<h3 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">Beautiful Reading UI</h3>
					<p class="text-gray-700 dark:text-gray-300 leading-relaxed">
						Enjoy an immersive reading experience with our carefully crafted interface designed for comfort and focus.
					</p>
				</div>
				
				<!-- Feature 2 -->
				<div class="acrylic rounded-2xl p-8 hover-lift  ">
					<div class="w-16 h-16 bg-gradient-to-br from-green-500 to-teal-600 rounded-2xl flex items-center justify-center mb-6">
						<Upload class="w-8 h-8 text-white" />
					</div>
					<h3 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">Upload Your Books</h3>
					<p class="text-gray-700 dark:text-gray-300 leading-relaxed">
						Import your own book collection and keep everything organized in one place. Support for multiple formats.
					</p>
				</div>
				
				<!-- Feature 3 -->
				<div class="acrylic rounded-2xl p-8 hover-lift  ">
					<div class="w-16 h-16 bg-gradient-to-br from-orange-500 to-red-600 rounded-2xl flex items-center justify-center mb-6">
						<Search class="w-8 h-8 text-white" />
					</div>
					<h3 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">Discover New Books</h3>
					<p class="text-gray-700 dark:text-gray-300 leading-relaxed">
						Explore our curated collection of free books from authors worldwide. Find your next favorite read.
					</p>
				</div>
			</div>
		</div>
	</section>
	
	<!-- Free Library Section -->
	<section class="py-20 px-4 sm:px-6 lg:px-8">
		<div class="max-w-4xl mx-auto">
			<div class="acrylic rounded-3xl p-8 lg:p-12 text-center">
				<div class=" ">
					<div class="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-yellow-400 to-orange-500 rounded-full mb-8">
						<Bookmark class="w-10 h-10 text-white" />
					</div>
					<h2 class="text- 3xl lg:text-5xl font-bold mb-6 text-gray-900 dark:text-white">Free Library Access</h2>
					<p class="text-xl text-gray-700 dark:text-gray-300 mb-8 max-w-2xl mx-auto">
						Access thousands of free books from generous authors. Build your digital library without any cost.
					</p>
					<div class="flex flex-col sm:flex-row gap-4 justify-center">
						<div class="flex items-center space-x-3">
							<div class="w-2 h-2 bg-green-400 rounded-full pulse-glow"></div>
							<span class="text-gray-700 dark:text-gray-300">Free to read</span>
						</div>
						<div class="flex items-center space-x-3">
							<div class="w-2 h-2 bg-blue-400 rounded-full pulse-glow"></div>
							<span class="text-gray-700 dark:text-gray-300">No subscription required</span>
						</div>
						<div class="flex items-center space-x-3">
							<div class="w-2 h-2 bg-purple-400 rounded-full pulse-glow"></div>
							<span class="text-gray-700 dark:text-gray-300">Offline reading</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>
	
	<!-- CTA Section -->
	<section class="py-20 px-4 sm:px-6 lg:px-8">
		<div class="max-w- 3xl mx-auto text-center  ">
			<h2 class="text- 3xl lg:text-5xl font-bold mb-6 text-gray-900 dark:text-white">Ready to Start Reading?</h2>
			<p class="text-xl text-gray-700 dark:text-gray-300 mb-8">
				Join thousands of readers who've made Books their digital home
			</p>
			<a href="/library" class="inline-block px-12 py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-full font-semibold text-xl hover:from-blue-600 hover:to-purple-700 transition-all duration-300 glow-effect hover-lift">
				Get Started Free
			</a>
		</div>
	</section>
</main>

<!-- Footer -->
<footer class="border-t border-gray-200 dark:border-gray-700 bg-white/50 dark:bg-black/50 backdrop-blur-sm">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
		<div class="flex flex-col md:flex-row justify-between items-center">
			<div class="flex items-center space-x-2 mb-4 md:mb-0">
				<div class="w-8 h-8 bg-gradient-to-br from-blue-400 to-purple-500 rounded-lg flex items-center justify-center">
					<Book class="w-5 h-5 text-white" />
				</div>
				<span class="text-xl font-bold text-gray-900 dark:text-white">Books</span>
			</div>
			<div class="flex space-x-6 text-gray-600 dark:text-gray-400">
				<a href="/" class="hover:text-gray-900 dark:hover:text-white transition-colors">Privacy</a>
				<a href="/" class="hover:text-gray-900 dark:hover:text-white transition-colors">Terms</a>
				<a href="/" class="hover:text-gray-900 dark:hover:text-white transition-colors">Support</a>
				<a href="/" class="hover:text-gray-900 dark:hover:text-white transition-colors">Contact</a>
			</div>
		</div>
		<div class="mt-8 pt-8 border-t border-gray-200 dark:border-gray-700 text-center">
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
		0%, 100% { transform: translateY(0px); }
		50% { transform: translateY(-10px); }
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

