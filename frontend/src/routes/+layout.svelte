<script lang="ts">
	import '../app.css';
	let { children } = $props();
	// No import needed here anymore!
</script>

<enhanced:img
	src="../lib/assets/bg.webp?w=480;768;1280;1920&format=avif;webp"
	alt="A dark, abstract background with a soft texture."
	sizes="100vw"
	class="fixed inset-0 -z-20 h-full w-full object-cover object-top"
/>

<div class="fixed inset-0 -z-10 bg-white/60 backdrop-blur-md dark:bg-black/40"></div>

<div class="relative flex min-h-screen flex-col">
	{@render children()}
</div>
