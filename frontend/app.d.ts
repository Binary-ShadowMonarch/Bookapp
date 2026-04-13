// src/app.d.ts
/// <reference types="@sveltejs/kit" />

declare namespace App {
	// shape of locals injected by hooks.server.ts
	interface Locals {
		user?: {
			email: string;
			// add any other fields you attach to locals.user
		};
	}

	// if you return `user` from load, you can type PageData too
	interface PageData {
		user: {
			email: string;
			// …and any other data you pass to the page
		};
	}
}
