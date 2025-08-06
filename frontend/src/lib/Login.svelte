<!-- this is the login page component -->
<!-- it handles both regular email/password login and Google OAuth -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { Book, Eye, EyeOff, Mail, Lock } from 'lucide-svelte';
	import { onMount } from 'svelte';

	// form data that gets updated as the user types
	let email = '';
	let password = '';
	let error: string | null = null; // error message to show to the user
	let showPassword = false; // toggle for showing/hiding password
	let isSubmitting = false; // prevents double-submission

	// This function runs when the component is first loaded
	onMount(async () => {
		// Check for an error message in the URL from the Google login redirect
		const urlParams = new URLSearchParams(window.location.search);
		const errorParam = urlParams.get('error');

		if (errorParam) {
			error = errorParam;
			// Clean up the URL so the error message doesn't stay in the address bar
			window.history.replaceState({}, document.title, '/login');
		}

		// Then, check if the user is already logged in
		await check();
	});

	// check if the user is already logged in when the page loads
	// if they are, redirect them to the library
	async function check() {
		// console.log('DEBUG: Checking if user is already logged in');
		try {
			let res = await fetch('/api/protected/profile', {
				method: 'GET',
				credentials: 'include'
			});

			if (res.ok) {
				// console.log('DEBUG: User is already logged in, redirecting to library');
				goto('/library');
				return;
			}
			// if the request failed with 401, try to refresh the token
			if (res.status === 401) {
				// console.log('DEBUG: Token expired, attempting refresh');
				const refresh = await fetch('/api/refresh', {
					method: 'POST',
					credentials: 'include'
				});
				if (refresh.ok) {
					// console.log('DEBUG: Token refreshed successfully, redirecting to library');
					goto('/library');
					return;
				}
			}
		} catch (e) {
			console.error('DEBUG: Failed to check login status:', e);
		}
	}

	// form validation - these are reactive variables that update as the user types
	$: isValidEmail = email.includes('@') && email.includes('.'); // basic email validation
	$: isValidPassword = password.length >= 8; // password must be at least 8 characters
	$: isFormValid = isValidEmail && isValidPassword && !isSubmitting; // form is valid if all conditions are met

	// toggle password visibility
	function togglePassword() {
		// console.log('DEBUG: Toggling password visibility');
		showPassword = !showPassword;
	}

	// handle the form submission when user clicks "Sign In"
	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		error = null;
		isSubmitting = true;

		// console.log('DEBUG: Attempting login for:', email);

		try {
			// send the login request to my backend
			const res = await fetch('/api/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
				body: new URLSearchParams({ mail: email, password }),
				credentials: 'include' // this sends the cookies
			});

			if (!res.ok) {
				// try to get a meaningful error message from the response
				let message = 'Login failed';
				const contentType = res.headers.get('content-type') || '';

				if (contentType.includes('application/json')) {
					const json = await res.json();
					if (json.message) message = json.message;
				} else {
					message = (await res.text()).trim();
				}
				throw new Error(message);
			}

			// login was successful! redirect to the library
			// console.log('DEBUG: Login successful, redirecting to library');
			goto('/library');
		} catch (err: any) {
			// console.log('DEBUG: Login failed:', err.message);
			// handle different types of errors and show appropriate messages
			if (err.message === 'invalid credentials') {
				error = 'Invalid credentials. Please check your email and password.';
			} else if (err.message === 'user not found') {
				error = 'User not found. Please check your email or sign up.';
			} else if (err.message === 'ACCOUNT ASSOCIATED WITH THIS EMAIL EXISTS SIGN IN USING GOOGLE') {
				error = 'An account with this email already exists. Please sign in using Google.';
			} else {
				// for unexpected errors, show a generic message
				error = 'An unexpected error occurred. Please try again later.';
			}
		} finally {
			isSubmitting = false;
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
	<div class="form-group">
		<label for="mail" class="form-label">
			<Mail class="h-4 w-4" />
			Email Address
		</label>
		<input
			id="mail"
			name="mail"
			type="email"
			bind:value={email}
			placeholder="your@email.com"
			class="form-input"
			class:error={email && !isValidEmail}
			required
		/>
		{#if email && !isValidEmail}
			<span class="field-error">Please enter a valid email address</span>
		{/if}
	</div>

	<div class="form-group">
		<label for="password" class="form-label">
			<Lock class="h-4 w-4" />
			Password
		</label>
		<div class="password-input-container">
			<input
				id="password"
				name="password"
				type={showPassword ? 'text' : 'password'}
				bind:value={password}
				placeholder="Enter your password"
				class="form-input password-input"
				class:error={password && !isValidPassword}
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
			</button>
		</div>
		{#if password && !isValidPassword}
			<span class="field-error">Password must be at least 8 characters long</span>
		{/if}
	</div>

	<!-- Centered Forgot Password as button for accessibility -->
	<div class="forgot-container">
		<button type="button" class="link">Forgot password</button>
	</div>
	{#if error}
		<p class="text-center font-semibold text-red-400">{error}</p>
	{/if}
	<button type="submit" class="submit-button" disabled={!isFormValid}>
		{isSubmitting ? 'Signing In...' : 'Sign In'}
	</button>
	<div class="flex-col">
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
	</div>
</form>

<style>
	.form {
		max-width: 430px;
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
		background-color: #1f1f1f;
		padding: 2rem;
		width: 28rem;
		border-radius: 1rem;
		font-family:
			-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans',
			'Helvetica Neue', sans-serif;
		background-color: rgba(0, 0, 0, 0.2);
		backdrop-filter: blur(50px);
	}

	.form-group {
		margin-bottom: 1.5rem;
	}

	.form-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-bottom: 0.5rem;
		font-size: 0.9rem;
		font-weight: 500;
		color: rgba(255, 255, 255, 0.9);
	}

	.form-input {
		width: 100%;
		padding: 0.75rem;
		font-size: 1rem;
		border: 1px solid rgba(255, 255, 255, 0.3);
		border-radius: 0.5rem;
		background: rgba(0, 0, 0, 0.2);
		backdrop-filter: blur(10px);
		color: white;
		transition: all 0.2s ease;
		box-sizing: border-box;
		/* margin-bottom: 0%; */
	}

	.form-input:focus {
		outline: none;
		border-color: #d0b3ff;
		box-shadow: 0 0 0 2px rgba(208, 179, 255, 0.2);
	}

	.field-error {
		display: block;
		margin-top: 0.25rem;
		font-size: 0.8rem;
		color: #ff6b6b;
	}
	.form-input.error {
		border-color: #ff6b6b;
	}

	.form-input::placeholder {
		color: rgba(255, 255, 255, 0.5);
	}

	.password-input-container {
		position: relative;
	}

	.password-input {
		padding-right: 2.5rem;
	}

	.password-toggle {
		position: absolute;
		right: 0.75rem;
		top: 50%;
		transform: translateY(-50%);
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.6);
		cursor: pointer;
		transition: color 0.2s ease;
	}

	.password-toggle:hover {
		color: rgba(255, 255, 255, 0.9);
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

	.submit-button {
		width: 100%;
		padding: 0.8rem;
		margin-top: 0.5rem;
		font-size: 1rem;
		font-weight: 600;
		border: none;
		border-radius: 0.5rem;
		background: linear-gradient(145deg, #181822, #600c85);
		color: white;
		cursor: pointer;
		transition: all 0.2s ease;
		margin-bottom: 0.5rem;
	}

	.submit-button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none;
	}

	.submit-button:not(:disabled):hover {
		transform: scale(1.02);
		box-shadow: 0 4px 15px rgba(0, 0, 0, 0.4);
	}

	.link:focus,
	.link-inline:focus {
		outline: 2px solid #fff;
		outline-offset: 2px;
	}

	.p {
		text-align: center;
		color: #f1f1f1;
		font-size: 14px;
		margin: 0rem;
	}

	.p.line {
		margin-top: 10px;
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

	.password-input-container {
		position: relative;
	}

	.password-input {
		padding-right: 2.5rem;
	}

	.password-toggle {
		position: absolute;
		right: 0.75rem;
		top: 50%;
		transform: translateY(-50%);
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.6);
		cursor: pointer;
		transition: color 0.2s ease;
	}

	.password-toggle:hover {
		color: rgba(255, 255, 255, 0.9);
	}

	.btn.google:hover
    /* ,.btn.github:hover  */ {
		border: 1px solid #2d79f3;
	}

	.text-red-400 {
		color: #ff6b6b; /* A more vibrant red for errors */
	}
</style>
