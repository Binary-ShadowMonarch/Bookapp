// src/routes/signup/+page.server.ts

import type { Actions } from './$types';
import { fail, redirect } from '@sveltejs/kit';

export const actions: Actions = {
  default: async ({ request, fetch, cookies }) => {
    const form = await request.formData();
    const mail = form.get('mail')?.toString() ?? '';
    const password = form.get('password')?.toString() ?? '';

    // 1) Call Go backend
    const res = await fetch('http://localhost:8080/register', {
      method: 'POST',
      body: new URLSearchParams({ mail, password })
    });

    if (!res.ok) {
      // If registration failed, re‑render form with error
      const error = await res.text();
      return fail(res.status, { mail, error });
    }

    // 2) On success, set one‑time cookie
    cookies.set('registered', 'true', {
      httpOnly: true,
      maxAge: 60,    // valid for 60 seconds
      path: '/'      // required
    });

    // 3) Redirect to success page
    throw redirect(303, '/signup/success');
  }
};
