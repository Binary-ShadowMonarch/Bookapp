<script lang="ts">
  import { Book } from 'lucide-svelte';
  /**
   * Controls the visibility of the secondary, hidden-by-default navigation row.
   * @type {boolean}
   */
  export let showSecondaryRow: boolean = false;
</script>

<nav class="sticky top-0 left-0 z-30 w-full bg-white/80 shadow backdrop-blur-md dark:border-black/30 dark:bg-black/60">
  <div class="flex h-14 items-center px-4 md:px-6">
    <div class="flex flex-1 justify-left gap-2">
      <div class="text-lg font-bold text-gray-900 dark:text-gray-100 select-none">
        <a href="/" class="flex items-center gap-2" aria-label="Homepage">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-400 to-purple-500">
            <Book class="h-5 w-5 text-white" />
          </div>
          <slot name="name"></slot>
        </a>
      </div>
    </div>

    <div class=" md:flex flex-1 items-center justify-center">
      <slot name="center" />
    </div>

    <div class="flex flex-1 items-center justify-end gap-2">
      <slot name="right-buttons" />
    </div>
  </div>

  {#if showSecondaryRow}
    <div class="flex h-14 items-center justify-center px-4 md:px-6">
      <slot name="secondary-center-buttons" />
    </div>
  {/if}
</nav>