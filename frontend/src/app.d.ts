/// <reference types="@sveltejs/kit" />

// See https://kit.svelte.dev/docs/types#the-app-namespace
// for information about these interfaces
declare namespace App {
	interface Locals {
		user: {
			isError: boolean;
			isExists: boolean;
			isAdmin: boolean;
			username: string;
			token: string | null;
		};
	}
	// interface Platform {}
	interface Session {
		user: {
			isError: boolean;
			isExists: boolean;
			isAdmin: boolean;
			username: string;
			token: string | null;
		};
	}
	interface Stuff {
		title?: string;
		description?: string;
		author?: string;
	}
}
