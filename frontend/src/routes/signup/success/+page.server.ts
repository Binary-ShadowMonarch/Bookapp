import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = ({ cookies }) => {
	// Only allow access if user just completed registration
	if (!cookies.get('registered')) {
		throw redirect(303, '/signup');
	}

	// Clear the registration cookie after successful access
	cookies.delete('registered', { path: '/' });

	return {};
};
