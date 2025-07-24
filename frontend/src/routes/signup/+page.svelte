<!-- src/routes/signup/+page.svelte -->
<script lang="ts">
	import Signup from '$lib/Signup.svelte';
	import type { ActionData } from './$types';

	export let form: ActionData;
</script>

<svelte:head>
	<title>Sign Up - Books</title>
</svelte:head>
<div class="flex min-h-screen items-center justify-center">
	<Signup {form} />
</div>
