// src/hooks.server.ts
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    // Hit /me to check if the user is logged in
    let res = await event.fetch('http://localhost:8080/library');

    if (res.status === 401) {
        // Try refreshing tokens
        const refreshRes = await event.fetch('http://localhost:8080/refresh', {
            method: 'POST'
        });

        if (refreshRes.ok) {
            // Retry the /me endpoint
            res = await event.fetch('http://localhost:8080/library');
        } else {
            // Cleanup cookies (optional)
            event.cookies.delete('access_token', { path: '/' });
            event.cookies.delete('refresh_token', { path: '/refresh' });
        }
    }

    return resolve(event);
};
