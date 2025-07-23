// src/routes/signup/verify/+page.server.ts
import type { Actions, PageServerLoad } from './$types';
import { fail, redirect } from '@sveltejs/kit';

export const load: PageServerLoad = ({ cookies, url }) => {
    // Check if user has pending registration
    if (!cookies.get('pending')) {
        throw redirect(303, '/signup');
    }


    // Get email from URL or cookie
    const mail = url.searchParams.get('mail') || '';
    if (!mail) {
        throw redirect(303, '/signup');
    }
    cookies.delete("pending", { path: '/' })
    return {
        mail
    };
};

export const actions: Actions = {
    // Verify action for verifying the code

    verify: async ({ request, fetch, cookies }) => {
        const data = await request.formData();
        const mail = data.get('mail')?.toString() || '';
        const code = data.get('code')?.toString() || '';

        if (!mail || !code) {
            return fail(400, { mail, error: 'Email and verification code are required.' });
        }

        try {
            const res = await fetch('http://localhost:8080/signup/verify', {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({ mail, code })
            });

            if (!res.ok) {
                const err = await res.text();
                return fail(res.status, { mail, error: err });
            }

            // Clear pending cookie and set registered cookie
            cookies.delete('pending', { path: '/' });
            cookies.set('registered', 'true', {
                path: '/',
                httpOnly: true,
                secure: true,
                maxAge: 60, // 1 minute to access success page
                sameSite: 'lax'
            });
        } catch (error) {
            return fail(500, { mail, error: 'Verification failed. Please try again.' });
        }

        // Success - redirect to success page
        throw redirect(303, '/signup/success');
    },

    // Resend action for sending a new code
    resend: async ({ request, fetch }) => {
        const data = await request.formData();
        const mail = data.get('mail')?.toString() || '';

        if (!mail) {
            return fail(400, { mail, error: 'Email is required to resend code.' });
        }

        try {
            console.log(mail)
            const res = await fetch('http://localhost:8080/signup/resend', {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({ mail })
            });

            if (!res.ok) {
                const err = await res.text();
                return fail(res.status, { mail, error: err });
            }

            return { mail, resendSuccess: true };
        } catch (error) {
            return fail(500, { mail, error: 'Failed to resend code. Please try again.' });
        }
    }
};