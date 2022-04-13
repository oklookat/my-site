/// <reference types="@sveltejs/kit" />

// TODO: after updating vite or svelte little broken and need fix like this
// https://github.com/sveltejs/kit/discussions/4551#discussioncomment-2532687
// (13.04.2022)
// as soon as possible, you need to update again, suddenly the error disappears

// See https://kit.svelte.dev/docs/types#the-app-namespace
// for information about these interfaces
declare namespace App {
	interface Locals {
		user: {
			isExists: boolean
			isAdmin: boolean
			username: string
			token: string | null
		}
	}
	// interface Platform {}
	interface Session {
	}
	// interface Stuff {}
}
