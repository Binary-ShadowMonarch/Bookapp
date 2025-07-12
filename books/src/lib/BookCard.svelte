<script lang="ts">
  export let title: string;
  export let author: string;
  export let image: string = '/default.webp';
  export let status: 'read' | 'unread' | 'finished' = 'unread';
  export let completion: number = 0;

  function handleImgError(event: Event) {
    const target = event.target as HTMLImageElement | null;
    if (target) target.src = '/default.webp';
  }
</script>

<div class="group relative flex flex-col justify-end aspect-[2/3] rounded-2xl overflow-hidden shadow-lg bg-white/60 dark:bg-black/40 backdrop-blur-md transition-transform hover:scale-105 cursor-pointer min-w-[140px] max-w-[180px] sm:min-w-[160px] sm:max-w-[220px] md:min-w-[180px] md:max-w-[240px] lg:min-w-[200px] lg:max-w-[260px]">
  <img
    src={image}
    alt={title}
    class="absolute inset-0 w-full h-full object-cover object-center z-0 group-hover:brightness-90 transition"
    on:error={handleImgError}
  />
  <!-- Info background -->
  <div class="relative z-10 p-3 pb-4 bg-white/90 dark:bg-black/80 rounded-b-2xl w-full mt-auto shadow-md flex flex-col">
    <h3 class="text-base font-semibold truncate text-gray-900 dark:text-gray-100">{title}</h3>
    <p class="text-sm text-gray-700 dark:text-gray-300 truncate">{author}</p>
    <span class="inline-flex items-center gap-2 mt-2 px-2 py-0.5 rounded-full text-xs font-medium shadow-sm border
      {status === 'read' ? 'bg-green-200 text-green-800 dark:bg-green-700 dark:text-green-100 border-green-300 dark:border-green-600' : ''}
      {status === 'unread' ? 'bg-yellow-200 text-yellow-800 dark:bg-yellow-700 dark:text-yellow-100 border-yellow-300 dark:border-yellow-600' : ''}
      {status === 'finished' ? 'bg-blue-200 text-blue-800 dark:bg-blue-700 dark:text-blue-100 border-blue-300 dark:border-blue-600' : ''}">
      {status.charAt(0).toUpperCase() + status.slice(1)}
      <span class="block w-12 h-1 rounded-full bg-gray-300 dark:bg-gray-700 overflow-hidden">
        <span class="block h-1 rounded-full transition-all duration-300"
          class:bg-green-500={completion === 100}
          class:bg-blue-500={completion < 100 && completion > 0}
          class:bg-gray-400={completion === 0}
          style="width: {completion}%"></span>
      </span>
      <span class="ml-1 font-semibold">{completion}%</span>
    </span>
  </div>
</div> 