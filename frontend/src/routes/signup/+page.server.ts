// src/routes/signup/+page.server.ts
import type { Actions, PageServerLoad } from './$types';
import { fail, redirect } from '@sveltejs/kit';

export const load: PageServerLoad = ({ cookies }) => {
    // Redirect if already registered
    if (cookies.get('registered')) {
        throw redirect(303, '/');
    }
    // Clear any existing pending registration
    if (cookies.get('pending')) {
        cookies.delete('pending', { path: '/' });
    }
};

export const actions: Actions = {
    default: async ({ request, fetch, cookies }) => {
        const data = await request.formData();
        const mail = data.get('mail')?.toString() || '';
        const password = data.get('password')?.toString() || '';

        if (!mail || !password) {
            return fail(400, { mail, error: 'Email and password are required.' });
        }

        // Basic validation
        if (password.length < 8) {
            return fail(400, { mail, error: 'Password must be at least 8 characters long.' });
        }
        if (!mail.includes('@') || !mail.includes('.')) {
            return fail(400, { mail, error: 'Please enter a valid email address.' });
        }

        try {
            const res = await fetch('http://backend:8080/api/register/request', {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({ mail, password })
            });

            if (!res.ok) {
                const err = await res.text();
                console.log(err)
                return fail(res.status, { mail, error: err.toUpperCase() });
            }
            const token = crypto.randomUUID();
            // Set pending cookie for verification step
            cookies.set('pending', token, {
                path: '/',
                httpOnly: true,
                secure: true,
                maxAge: 300, // 5 minutes to match verification expiry
                sameSite: 'lax'
            });

        } catch (error) {
            // Handle network/server errors only (not redirects)
            return fail(500, { mail, error: 'Registration failed. Please try again.' });
        }

        // Success - redirect to verification page
        throw redirect(303, `/signup/verify?mail=${encodeURIComponent(mail)}`);
    }
};