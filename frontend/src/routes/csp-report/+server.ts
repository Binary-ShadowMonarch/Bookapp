// src/routes/csp-report/+server.ts

import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

export const POST: RequestHandler = async (event) => {
	try {
		// Parse the incoming CSP report payload
		const report = await event.request.json();

		// Log the CSP violation (in production, send to a monitoring service)
		console.warn('CSP Violation:', {
			timestamp: new Date().toISOString(),
			userAgent: event.request.headers.get('user-agent'),
			report
		});

		// Optionally store the violation in your database or monitoring system
		// await storeViolation(report);

		return json({ received: true });
	} catch (error) {
		console.error('Error processing CSP report:', error);
		return json({ error: 'Failed to process report' }, { status: 500 });
	}
};
