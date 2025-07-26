<!-- src/lib/Login.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { Book, Eye, EyeOff } from 'lucide-svelte';
	import { onMount } from 'svelte';
	
	let email = '';
	let password = '';
	let error: string | null = null;
	let showPassword = false;

	async function check() {
		let res = await fetch('/api/protected/profile', {
			method: 'GET',
			credentials: 'include'
		});

		if (res.ok) {
			goto('/library');
			return;
		}
		// 2) if unauthorized, try a refresh and retry
		if (res.status === 401) {
			const refresh = await fetch('/api/refresh', {
				method: 'POST',
				credentials: 'include'
			});
			if (refresh.ok) {
				// retry the library fetch
				res = await fetch('/api/protected/library', {
					method: 'GET',
					credentials: 'include'
				});
				goto('/library');
				return;
			}
		}
	}
	
	onMount(async () => {
		check();
	});

	function togglePassword() {
		showPassword = !showPassword;
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		error = null;

		try {
			const res = await fetch('/api/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
				body: new URLSearchParams({ mail: email, password }),
				credentials: 'include'
			});

			if (!res.ok) {
				let message = 'Login failed';
				const contentType = res.headers.get('content-type') || '';

				if (contentType.includes('application/json')) {
					const json = await res.json();
					if (json.message) message = json.message;
				} else {
					message = await res.text();
				}
				throw new Error(message);
			}

			// Redirect after successful login. SvelteKit's layout.server.ts will now
			// correctly identify the user as authenticated based on new cookies.
			console.log('Logged in successfully!');
			goto('/library');
		} catch (err: any) {
			error = err.message || 'Unexpected error occurred';
		}
	}
</script>

<form class="form" on:submit|preventDefault={handleSubmit}>
	<div class="flex select-none justify-center text-lg font-bold text-gray-900 dark:text-gray-100">
		<a href="/" class="flex items-center gap-2" aria-label="Homepage">
			<div
				class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-400 to-purple-500"
			>
				<Book class="h-5 w-5 text-white" />
			</div>
			<h1>Books</h1>
		</a>
	</div>
	<div class="flex-column">
		<label for="email">Email</label>
	</div>
	<div class="inputForm">
		<!-- Email icon -->
		<svg height="20" viewBox="0 0 32 32" width="20" xmlns="http://www.w3.org/2000/svg">
			<g id="Layer_3" data-name="Layer 3">
				<path
					d="m30.853 13.87a15 15 0 0 0 -29.729 4.082 15.1 15.1 0 0 0 12.876 12.918 15.6 15.6 0 0 0 2.016.13 14.85 14.85 0 0 0 7.715-2.145 1 1 0 1 0 -1.031-1.711 13.007 13.007 0 1 1 5.458-6.529 2.149 2.149 0 0 1 -4.158-.759v-10.856a1 1 0 0 0 -2 0v1.726a8 8 0 1 0 .2 10.325 4.135 4.135 0 0 0 7.83.274 15.2 15.2 0 0 0 .823-7.455zm-14.853 8.13a6 6 0 1 1 6-6 6.006 6.006 0 0 1 -6 6z"
				/>
			</g>
		</svg>
		<input
			id="email"
			type="email"
			class="input"
			placeholder="Enter your Email"
			bind:value={email}
			required
		/>
	</div>

	<div class="flex-column">
		<label for="password">Password</label>
	</div>
	<div class="inputForm password-container">
		<!-- Password icon -->
		<svg height="20" viewBox="-64 0 512 512" width="20" xmlns="http://www.w3.org/2000/svg">
			<path
				d="m336 512h-288c-26.453125 0-48-21.523438-48-48v-224c0-26.476562 21.546875-48 48-48h288c26.453125 0 48 21.523438 48 48v224c0 26.476562-21.546875 48-48 48zm-288-288c-8.8125 0-16 7.167969-16 16v224c0 8.832031 7.1875 16 16 16h288c8.8125 0 16-7.167969 16-16v-224c0-8.832031-7.1875-16-16-16zm0 0"
			/>
			<path
				d="m304 224c-8.832031 0-16-7.167969-16-16v-80c0-52.929688-43.070312-96-96-96s-96 43.070312-96 96v80c0 8.832031-7.167969 16-16 16s-16-7.167969-16-16v-80c0-70.59375 57.40625-128 128-128s128 57.40625 128 128v80c0 8.832031-7.167969 16-16 16zm0 0"
			/>
		</svg>
							<div class="password-input-container">

		<input
			id="password"
			type={showPassword ? 'text' : 'password'}
			class="input password-input"
			placeholder="Enter your Password"
			bind:value={password}
			required
		/>
						<button
							type="button"
							class="password-toggle"
							on:click={togglePassword}
							aria-label={showPassword ? 'Hide password' : 'Show password'}
						>
							{#if showPassword}
								<EyeOff class="h-4 w-4" />
							{:else}
								<Eye class="h-4 w-4" />
							{/if}
	</div>
</div>

	<!-- Centered Forgot Password as button for accessibility -->
	<div class="forgot-container">
		<button type="button" class="link">Forgot password?</button>
	</div>
	{#if error}
		<p class="text-red-400">{error}</p>
	{/if}
	<button type="submit" class="button-submit signin">Sign In</button>

	<p class="p">
		Don't have an account? <button type="button" class="link-inline"
			><a href="/signup">Sign Up</a></button
		>
	</p>

	<p class="p line">Or Continue Using</p>

	<button type="button" class="btn google">
		<a class="flex gap-2" aria-label="GoogleAuth" href="/api/auth/google/login">
			<svg
				version="1.1"
				width="20"
				id="Layer_1"
				xmlns="http://www.w3.org/2000/svg"
				xmlns:xlink="http://www.w3.org/1999/xlink"
				x="0px"
				y="0px"
				viewBox="0 0 512 512"
				style="enable-background:new 0 0 512 512;"
				xml:space="preserve"
			>
				<path
					style="fill:#FBBB00;"
					d="M113.47,309.408L95.648,375.94l-65.139,1.378C11.042,341.211,0,299.9,0,256
		c0-42.451,10.324-82.483,28.624-117.732h0.014l57.992,10.632l25.404,57.644c-5.317,15.501-8.215,32.141-8.215,49.456
		C103.821,274.792,107.225,292.797,113.47,309.408z"
				></path>
				<path
					style="fill:#518EF8;"
					d="M507.527,208.176C510.467,223.662,512,239.655,512,256c0,18.328-1.927,36.206-5.598,53.451
		c-12.462,58.683-45.025,109.925-90.134,146.187l-0.014-0.014l-73.044-3.727l-10.338-64.535
		c29.932-17.554,53.324-45.025,65.646-77.911h-136.89V208.176h138.887L507.527,208.176L507.527,208.176z"
				></path>
				<path
					style="fill:#28B446;"
					d="M416.253,455.624l0.014,0.014C372.396,490.901,316.666,512,256,512
		c-97.491,0-182.252-54.491-225.491-134.681l82.961-67.91c21.619,57.698,77.278,98.771,142.53,98.771
		c28.047,0,54.323-7.582,76.87-20.818L416.253,455.624z"
				></path>
				<path
					style="fill:#F14336;"
					d="M419.404,58.936l-82.933,67.896c-23.335-14.586-50.919-23.012-80.471-23.012
		c-66.729,0-123.429,42.957-143.965,102.724l-83.397-68.276h-0.014C71.23,56.123,157.06,0,256,0
		C318.115,0,375.068,22.126,419.404,58.936z"
			></path>
			</svg>
			<p>Google</p>
		</a>
	</button>
</form>

<style>
	.form {
		display: flex;
		flex-direction: column;
		gap: 10px;
		background-color: #1f1f1f;
		padding: 30px;
		width: 450px;
		border-radius: 20px;
		font-family:
			-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans',
			'Helvetica Neue', sans-serif;
		background-color: rgba(0, 0, 0, 0.2);
		backdrop-filter: blur(50px);
	}

	.flex-column > label {
		color: #f1f1f1;
		font-weight: 600;
	}

	.inputForm {
		border: 1.5px solid #333;
		border-radius: 10px;
		height: 50px;
		display: flex;
		align-items: center;
		padding-left: 10px;
		transition: 0.2s ease-in-out;
		background-color: #2b2b2b;
	}

	.password-container {
		position: relative;
	}

	.input {
		margin-left: 10px;
		border-radius: 10px;
		border: none;
		width: 100%;
		height: 100%;
		background-color: #2b2b2b;
		color: #f1f1f1;
	}

	.password-input {
		padding-right: 40px;
	}

	.input:focus {
		outline: none;
	}

	.password-toggle {
		position: absolute;
		right: 10px;
		top: 50%;
		transform: translateY(-50%);
		background: none;
		border: none;
		color: rgba(241, 241, 241, 0.6);
		cursor: pointer;
		transition: color 0.2s ease;
		display: flex;
		align-items: center;
		padding: 0;
	}

	.password-toggle:hover {
		color: rgba(241, 241, 241, 0.9);
	}

	.forgot-container {
		display: flex;
		justify-content: center;
	}

	/* Accessible link-style buttons */
	.link,
	.link-inline {
		background: none;
		border: none;
		color: #2d79f3;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		text-decoration: underline;
		padding: 0;
	}

	.button-submit {
		margin: 20px 0 10px 0;
		background-color: #2d79f3;
		border: none;
		color: white;
		font-size: 15px;
		font-weight: 500;
		border-radius: 10px;
		height: 50px;
		width: 100%;
		cursor: pointer;
	}

	.button-submit:focus,
	.link:focus,
	.link-inline:focus {
		outline: 2px solid #fff;
		outline-offset: 2px;
	}

	.p {
		text-align: center;
		color: #f1f1f1;
		font-size: 14px;
		margin: 5px 0;
	}

	.p.line {
		margin-top: 10px;
	}

	.signin {
		width: 100%;
		padding: 0.75rem;
		font-size: 0.9rem;
		font-weight: bold;
		border: none;
		border-radius: 0.5rem;
		background: linear-gradient(145deg, #181822, #600c85);
		color: white;
		cursor: pointer;
		transition:
			transform 0.2s,
			box-shadow 0.2s;
	}
	.signin:hover {
		transform: scale(1.02);
		box-shadow: 0 4px 15px rgba(0, 0, 0, 0.4);
	}

	.btn {
		margin-top: 10px;
		width: 100%;
		height: 50px;
		border-radius: 10px;
		display: flex;
		justify-content: center;
		align-items: center;
		font-weight: 500;
		gap: 10px;
		border: 1px solid #333;
		background-color: #2b2b2b;
		color: #f1f1f1;
		cursor: pointer;
		transition: 0.2s ease-in-out;
	}

	.btn.google:hover,
	.btn.github:hover {
		border: 1px solid #2d79f3;
	}
	</style>