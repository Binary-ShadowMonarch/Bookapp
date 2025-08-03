<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<!-- Background image with acrylic/blur effect -->
<div class="fixed inset-0 -z-10 bg-cover bg-center" style="background-image: url('/bg.webp');">
	<div class="h-full w-full bg-white/60 backdrop-blur-md dark:bg-black/40"></div>
	<!-- Main content -->
</div>
<div class="relative flex min-h-screen flex-col">
	{@render children()}
</div>
