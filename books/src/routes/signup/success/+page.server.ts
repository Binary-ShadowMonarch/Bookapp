// src/routes/signup/success/+page.server.ts

import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ cookies }) => {
  // Check for our one‑time flag
  const ok = cookies.get('registered') === 'true';
  if (!ok) {
    // No flag → bounce back to signup
    throw redirect(303, '/signup');
  }

  // Clear the cookie so refreshing this page won’t work
  cookies.set('registered', '', {
    httpOnly: true,
    maxAge: -1,
    path: '/'   // required
  });
};
