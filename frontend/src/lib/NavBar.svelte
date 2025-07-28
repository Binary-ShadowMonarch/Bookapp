<script lang="ts">
	import { Book } from 'lucide-svelte';
	/**
	 * Controls the visibility of the secondary, hidden-by-default navigation row.
	 * @type {boolean}
	 */
	export let showSecondaryRow: boolean = false;
</script>

<nav
	class="sticky left-0 top-0 z-30 w-full bg-white/80 shadow backdrop-blur-md dark:border-black/30 dark:bg-black/60"
>
	<div class="flex h-14 items-center px-4 md:px-6">
		<div class="justify-left flex flex-1 gap-2">
			<div class="select-none text-lg font-bold text-gray-900 dark:text-gray-100">
				<a href="/" class="flex items-center gap-2" aria-label="Homepage">
					<div
						class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-400 to-purple-500"
					>
						<Book class="h-5 w-5 text-white" />
					</div>
					<slot name="name"></slot>
				</a>
			</div>
		</div>

		<div class=" flex-1 items-center justify-center md:flex">
			<slot name="center" />
		</div>

		<div class="flex flex-1 items-center justify-end gap-1">
			<slot name="right-buttons" />
		</div>
	</div>

	{#if showSecondaryRow}
		<div class="flex h-14 items-center justify-center px-4 md:px-6">
			<slot name="secondary-center-buttons" />
		</div>
	{/if}
</nav>
