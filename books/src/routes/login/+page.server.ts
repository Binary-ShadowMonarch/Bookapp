// src/routes/login/+page.server.ts

import { fail, redirect } from '@sveltejs/kit';

const BASE_URL = 'http://localhost:8080';

export const actions = {
  default: async ({ request }) => {
    const data = await request.formData();
    const email = data.get('email') as string;
    const password = data.get('password') as string;

    if (!email || !password) {
      return fail(400, { error: 'Email and password are required.' });
    }

    const response = await fetch(`${BASE_URL}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({ mail: email, password: password }),
    });

    if (!response.ok) {
      return fail(401, { error: 'Invalid credentials.' });
    }

    // The Go backend sets the HttpOnly refresh token cookie. That's all we need.
    // We just redirect, and the root layout will handle fetching the access token.
    throw redirect(303, '/dashboard');
  },
};