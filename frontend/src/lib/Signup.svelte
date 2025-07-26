<!-- src/lib/Signup.svelte -->
<script lang="ts">
	import { Book, Eye, EyeOff, Mail, Lock } from 'lucide-svelte';
	import { enhance } from '$app/forms';
	import { onMount } from 'svelte';
	import type { ActionData } from '../routes/signup/$types';

	// Props - form data from server action
	export let form: ActionData | null = null;

	// Component state
	let isSubmitting = false;
	let showPassword = false;
	let mail = form?.mail || '';
	let password = '';
	let hasNavigated = false;

	// Reactive server error
	$: serverError = form?.error;

	// Form validation
	$: isValidEmail = mail.includes('@') && mail.includes('.');
	$: isValidPassword = password.length >= 8;
	$: isFormValid = isValidEmail && isValidPassword && !isSubmitting;

	// Page protection - warn before leaving
	onMount(() => {
		hasNavigated = false;

		const handleBeforeUnload = (e: BeforeUnloadEvent) => {
			if (mail || password) {
				e.preventDefault();

				return '';
			}
		};

		window.addEventListener('beforeunload', handleBeforeUnload);

		return () => {
			window.removeEventListener('beforeunload', handleBeforeUnload);
		};
	});

	function togglePassword() {
		showPassword = !showPassword;
	}
</script>

<div class="container">
	<div class="card">
		<div class="header">
			<a href="/" class="logo" aria-label="Homepage">
				<div class="logo-icon">
					<Book class="h-5 w-5 text-white" />
				</div>
				<h1>Books</h1>
			</a>
		</div>

		<div class="content">
			<div class="info-section">
				<h2>Create Your Account</h2>
				<p>Join our community and start building your digital library</p>
			</div>

			<form
				method="POST"
				use:enhance={() => {
					isSubmitting = true;
					hasNavigated = true;
					return async ({ update }) => {
						await update();
						isSubmitting = false;
					};
				}}
			>
				{#if serverError}
					<div class="server-error" role="alert">
						{serverError}
					</div>
				{/if}

				<div class="form-group">
					<label for="mail" class="form-label">
						<Mail class="h-4 w-4" />
						Email Address
					</label>
					<input
						id="mail"
						name="mail"
						type="email"
						bind:value={mail}
						placeholder="your@email.com"
						class="form-input"
						class:error={mail && !isValidEmail}
						autocomplete="email"
						required
					/>
					{#if mail && !isValidEmail}
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
							placeholder="Choose a strong password"
							class="form-input password-input"
							class:error={password && !isValidPassword}
							autocomplete="new-password"
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

				<button type="submit" class="submit-button" disabled={!isFormValid}>
					{isSubmitting ? 'Creating Account...' : 'Create Account'}
				</button>
			</form>

			<div class="footer justify-center">
				<p class="flex items-center justify-center gap-1">
					Already have an account? <a href="/login">Sign in</a> or
				</p>
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
		</div>
	</div>
</div>

<style>
	.container {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		padding: 2rem;
	}

	.card {
		width: 100%;
		max-width: 450px;
		background: rgba(0, 0, 0, 0.2);
		backdrop-filter: blur(50px);
		border: 1px solid rgba(255, 255, 255, 0.15);
		border-radius: 1rem;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
		overflow: hidden;
	}

	.header {
		padding: 1.5rem;
		text-align: center;
		border-bottom: 1px solid rgba(255, 255, 255, 0.2);
	}

	.logo {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		text-decoration: none;
		color: white;
		font-size: 1.125rem;
		font-weight: 700;
	}

	.logo-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2rem;
		height: 2rem;
		border-radius: 0.5rem;
		background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
	}

	.content {
		padding: 2rem;
	}

	.info-section {
		text-align: center;
		margin-bottom: 2rem;
	}

	.info-section h2 {
		margin: 0 0 0.5rem 0;
		font-size: 1.5rem;
		color: white;
	}

	.info-section p {
		margin: 0;
		color: rgba(255, 255, 255, 0.7);
		font-size: 0.9rem;
	}

	.server-error {
		color: #ffc9c9;
		background-color: rgba(255, 107, 107, 0.1);
		border: 1px solid rgba(255, 107, 107, 0.5);
		border-radius: 0.5rem;
		padding: 0.75rem;
		text-align: center;
		margin-bottom: 1.5rem;
		font-size: 0.9rem;
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
	}

	.form-input:focus {
		outline: none;
		border-color: #d0b3ff;
		box-shadow: 0 0 0 2px rgba(208, 179, 255, 0.2);
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

	.field-error {
		display: block;
		margin-top: 0.25rem;
		font-size: 0.8rem;
		color: #ff6b6b;
	}

	.submit-button {
		width: 100%;
		padding: 0.875rem;
		font-size: 1rem;
		font-weight: 600;
		border: none;
		border-radius: 0.5rem;
		background: linear-gradient(145deg, #181822, #600c85);
		color: white;
		cursor: pointer;
		transition: all 0.2s ease;
		margin-bottom: 1.5rem;
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

	.footer {
		text-align: center;
		padding-top: 0rem;
		padding:0.5rem;
		margin-bottom:0.5rem;
		border-top: 1px solid rgba(255, 255, 255, 0.1);
	}

	.footer p {
		/* margin: 0; */
		font-size: 0.9rem;
		color: rgba(255, 255, 255, 0.7);
	}

	.footer a {
		color: #d0b3ff;
		text-decoration: none;
		font-weight: 500;
		transition: color 0.2s ease;
	}

	.footer a:hover {
		color: #b084ff;
	}
	.btn {
		margin-top: 0px;
		padding-top: 0px;
		width: 100%;
		height: 50px;
		border-radius: 10px;
		display: flex;
		justify-content: center;
		align-items: center;
		font-weight: 500;
		gap: 0px;
		border: 1px solid #333;
		background-color: #2b2b2b;
		color: #f1f1f1;
		cursor: pointer;
		transition: 0.2s ease-in-out;
	}

	.btn.google:hover {
		border: 1px solid #2d79f3;
	}
</style>
