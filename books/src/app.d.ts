// src/app.d.ts

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: {
				isAuthenticated: boolean;
			}
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export { };