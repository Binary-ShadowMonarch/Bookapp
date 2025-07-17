// src/routes/library/+page.server.ts
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { goto } from '$app/navigation';

export const load: PageServerLoad = async ({ fetch }) => {
    // Check if the user is authenticated by calling your Go API
    const res = await fetch('http://localhost:8080/library');

    if (res.status === 401) {
        // Not authenticated → redirect to login
        throw redirect(303, '/login');
    }

    // goto("/library")
    // const user = await res.json();

    // Optionally, return user data to the page
    // return { user };
};
