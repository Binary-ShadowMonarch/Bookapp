<script lang="ts">
	import { Book } from 'lucide-svelte';

	// this controls whether to show the second row of navigation buttons
	// I use this for extra navigation options when needed
	export let showSecondaryRow: boolean = false;
</script>

<!-- this is the main navigation bar that stays at the top of the page -->
<!-- it has a nice blur effect and works in both light and dark mode -->
<nav
	class="sticky top-0 left-0 z-30 w-full bg-white/80 shadow backdrop-blur-md dark:border-black/30 dark:bg-black/60"
>
	<!-- main navigation row with logo, center content, and right buttons -->
	<div class="flex h-14 items-center px-4 md:px-6">
		<!-- left side - logo and app name -->
		<div class="justify-left flex flex-1 gap-2">
			<div class="text-lg font-bold text-gray-900 select-none dark:text-gray-100">
				<a href="/" class="flex items-center gap-2" aria-label="Homepage">
					<!-- this is my app logo - a book icon with a gradient background -->
					<div
						class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-400 to-purple-500"
					>
						<Book class="h-5 w-5 text-white" />
					</div>
					<!-- the app name goes here -->
					<slot name="name"></slot>
				</a>
			</div>
		</div>

		<!-- center section - I can put search or other navigation here -->
		<div class=" flex-1 items-center justify-center md:flex">
			<slot name="center" />
		</div>

		<!-- right side - login/logout buttons and other actions -->
		<div class="flex flex-1 items-center justify-end gap-1">
			<slot name="right-buttons" />
		</div>
	</div>

	<!-- optional second row for extra navigation buttons -->
	{#if showSecondaryRow}
		<div class="flex h-14 items-center justify-center px-4 md:px-6">
			<slot name="secondary-center-buttons" />
		</div>
	{/if}
</nav>
