<!-- src/routes/signup/verify/+page.svelte -->
<script lang="ts">
	import { Book } from 'lucide-svelte';
	import { enhance } from '$app/forms';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import type { PageData, ActionData } from './$types';

	// --- PROPS --
	export let data: PageData;
	export let form: ActionData;

	// --- REACTIVE STATE ---
	$: mail = form?.mail || data.mail;
	$: serverError = form?.error;
	$: resendSuccess = form?.resendSuccess;

	// --- LOCAL COMPONENT STATE ---
	let code = ['', '', '', '', '', ''];
	let codeInputs: HTMLInputElement[] = [];
	let expiryTime = 300; // 5 minutes
	let resendTime = 180; // 3 minutes for resend cooldown
	let canResend = false;
	let isVerifying = false;
	let isResending = false;

	$: isCodeComplete = code.every((digit) => digit !== '' && /^\d$/.test(digit));

	// --- LIFECYCLE & TIMERS ---
	onMount(() => {
		// Page protection - warn before leaving
		const handleBeforeUnload = (e: BeforeUnloadEvent) => {
			e.preventDefault();
			// e.returnValue = 'Changes in this site could be lost. Confirm resubmission?';

			return '';
		};
		window.addEventListener('beforeunload', handleBeforeUnload);

		// Expiry timer - redirects to signup when expired
		const expiryInterval = setInterval(() => {
			if (expiryTime > 0) {
				expiryTime--;
			} else {
				clearInterval(expiryInterval);
				goto('/signup');
			}
		}, 1000);

		// Resend timer - enables resend button after cooldown
		const resendInterval = setInterval(() => {
			if (resendTime > 0) {
				resendTime--;
			} else {
				canResend = true;
				clearInterval(resendInterval);
			}
		}, 1000);

		return () => {
			window.removeEventListener('beforeunload', handleBeforeUnload);
			clearInterval(expiryInterval);
			clearInterval(resendInterval);
		};
	});

	// Reset resend timer after successful resend
	$: if (resendSuccess) {
		resendTime = 180;
		canResend = false;
	}

	// --- EVENT HANDLERS ---
	function handleInput(index: number, event: Event) {
		const target = event.target as HTMLInputElement;
		const value = target.value;

		// Only allow single digit
		if (value.length > 1) {
			target.value = value.slice(-1);
		}

		code[index] = target.value;

		// Auto-focus next input
		if (target.value && index < 5) {
			codeInputs[index + 1]?.focus();
		}
	}

	function handleKeydown(index: number, event: KeyboardEvent) {
		// Handle backspace to move to previous input
		if (event.key === 'Backspace' && !code[index] && index > 0) {
			codeInputs[index - 1]?.focus();
		}

		// Handle paste
		if (event.key === 'v' && (event.ctrlKey || event.metaKey)) {
			event.preventDefault();
			navigator.clipboard.readText().then((text) => {
				const digits = text.replace(/\D/g, '').slice(0, 6);
				code = [...digits.padEnd(6, '').split('')];

				// Focus the next empty input or the last one
				const nextIndex = Math.min(digits.length, 5);
				codeInputs[nextIndex]?.focus();
			});
		}
	}

	function formatTime(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}
</script>

<svelte:head>
	<title>Verify Email - Books</title>
</svelte:head>

<div class="containerd">
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
				<h2>Verify Your Email</h2>
				<p>We've sent a verification code to <strong>{mail}</strong></p>
				<p class="expiry-info">
					This page will expire in <span class="time-highlight">{formatTime(expiryTime)}</span>
				</p>
			</div>

			<!-- Verification Form -->
			<form
				method="POST"
				action="?/verify"
				use:enhance={() => {
					isVerifying = true;
					return async ({ update }) => {
						await update();
						isVerifying = false;
					};
				}}
			>
				<input type="hidden" name="mail" value={mail} />
				<input type="hidden" name="code" value={code.join('')} />

				{#if serverError && !resendSuccess}
					<div class="server-error" role="alert">{serverError}</div>
				{/if}

				<div class="code-inputs">
					{#each code as _, index}
						<input
							bind:this={codeInputs[index]}
							bind:value={code[index]}
							type="text"
							inputmode="numeric"
							maxlength="1"
							on:input={(e) => handleInput(index, e)}
							on:keydown={(e) => handleKeydown(index, e)}
							class="code-input"
							autocomplete="one-time-code"
							aria-label={`Digit ${index + 1} of verification code`}
						/>
					{/each}
				</div>

				<button type="submit" class="verify-button" disabled={!isCodeComplete || isVerifying}>
					{isVerifying ? 'Verifying...' : 'Verify Account'}
				</button>
			</form>

			<!-- Resend Code Form -->
			<div class="resend-section">
				<form
					method="POST"
					action="?/resend"
					use:enhance={() => {
						isResending = true;
						return async ({ update }) => {
							await update();
							isResending = false;
						};
					}}
				>
					<input type="hidden" name="mail" value={mail} />
					<button type="submit" class="resend-button" disabled={!canResend || isResending}>
						{#if isResending}
							Sending...
						{:else if canResend}
							Resend Code
						{:else}
							Resend in {formatTime(resendTime)}
						{/if}
					</button>
				</form>
			</div>

			{#if resendSuccess}
				<div class="resend-success-message" role="status">A new code has been sent.</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.containerd {
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
		margin: 0 0 1rem 0;
		font-size: 1.5rem;
		color: white;
	}

	.info-section p {
		margin: 0.5rem 0;
		color: rgba(255, 255, 255, 0.8);
		font-size: 0.9rem;
	}

	.expiry-info {
		font-size: 0.8rem !important;
		color: rgba(255, 255, 255, 0.6) !important;
		margin-top: 1rem !important;
	}

	.time-highlight {
		color: #ff6b6b;
		font-weight: bold;
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

	.code-inputs {
		display: flex;
		justify-content: center;
		gap: 0.5rem;
		margin-bottom: 1.5rem;
	}

	.code-input {
		width: 3rem;
		height: 3rem;
		text-align: center;
		font-size: 1.5rem;
		font-weight: bold;
		border: 1px solid rgba(255, 255, 255, 0.3);
		border-radius: 0.5rem;
		background: rgba(0, 0, 0, 0.2);
		backdrop-filter: blur(10px);
		color: white;
		transition: all 0.2s ease;
	}

	.code-input:focus {
		outline: none;
		border-color: #d0b3ff;
		box-shadow: 0 0 0 2px rgba(208, 179, 255, 0.2);
		transform: scale(1.05);
	}

	.verify-button {
		width: 100%;
		padding: 0.75rem;
		font-size: 0.9rem;
		font-weight: bold;
		border: none;
		border-radius: 0.5rem;
		background: linear-gradient(145deg, #181822, #600c85);
		color: white;
		cursor: pointer;
		transition: all 0.2s ease;
		margin-bottom: 1rem;
	}

	.verify-button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none;
	}

	.verify-button:not(:disabled):hover {
		transform: scale(1.02);
		box-shadow: 0 4px 15px rgba(0, 0, 0, 0.4);
	}

	.resend-section {
		text-align: center;
	}

	.resend-button {
		background: transparent;
		border: 1px solid rgba(255, 255, 255, 0.3);
		color: rgba(255, 255, 255, 0.8);
		padding: 0.5rem 1rem;
		border-radius: 0.375rem;
		font-size: 0.8rem;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.resend-button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.resend-button:not(:disabled):hover {
		background: rgba(255, 255, 255, 0.1);
		border-color: #d0b3ff;
	}

	.resend-success-message {
		color: #bbf7d0;
		text-align: center;
		font-size: 0.85rem;
		margin-top: 1rem;
	}
</style>
