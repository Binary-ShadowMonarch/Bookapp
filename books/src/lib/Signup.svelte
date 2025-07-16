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
				e.returnValue = 'Changes in this site could be lost. Confirm resubmission?';
				return e.returnValue;
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

			<div class="footer">
				<p>Already have an account? <a href="/login">Sign in</a></p>
			</div>
		</div>
	</div>
</div>

<style>
	:global(body) {
		margin: 0;
		font-family: system-ui, sans-serif;
		color: white;
		background: #1a001f;
	}

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
		padding-top: 1rem;
		border-top: 1px solid rgba(255, 255, 255, 0.1);
	}

	.footer p {
		margin: 0;
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
</style>
