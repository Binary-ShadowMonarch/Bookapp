<!-- login +page.svelte -->
<script lang="ts">
	import LoginForm from '$lib/Login.svelte';
	import { goto } from '$app/navigation';
	// import { redirect } from '@sveltejs/kit';
</script>

<svelte:head>
	<title>Login - Books</title>
</svelte:head>

<div class=" relative flex min-h-screen items-center justify-center">
	<LoginForm />
</div>
