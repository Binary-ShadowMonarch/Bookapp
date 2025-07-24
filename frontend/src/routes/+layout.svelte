<script lang="ts">
	import '../app.css';
	let { children } = $props();
	
</script>

<!-- Background image with acrylic/blur effect -->
<div class="fixed inset-0 -z-10 bg-cover bg-center" style="background-image: url('/default.webp');">
  <div class="w-full h-full bg-white/60 dark:bg-black/40 backdrop-blur-md"></div>
  <!-- Main content -->
</div>
<div class="relative min-h-screen flex flex-col">
	{@render children()}
</div>

