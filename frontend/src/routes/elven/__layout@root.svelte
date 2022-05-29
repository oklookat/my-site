<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';

	/** check rights before access to admin panel */
	export const load: Load = async (event) => {
		const isError = event.session.user.isError;
		const isAdmin = event.session.user.isAdmin;
		const isLoginPage = isAdminPanelLoginPage(event.url);

		if (isError) {
			// redirect to main page if auth check error
			return {
				status: 302,
				redirect: '/'
			};
		}

		// redirect to login if not authorized
		if (!isAdmin) {
			// avoid redirect /login loop
			if (isLoginPage) {
				return {};
			}
			return {
				status: 302,
				redirect: '/elven/login'
			};
		}

		// redirect to main if admin on login page
		if (isLoginPage) {
			return {
				status: 302,
				redirect: '/elven'
			};
		}

		return {
			status: 200,
			props: {
				isAdmin: isAdmin
			}
		};
	};
</script>

<script lang="ts">
	import Header from '$elven/components/header.svelte';
	import ServiceWrapper from '$elven/components/service_wrapper.svelte';
	import Progress from '$elven/plugins/elvenProgress/progress.svelte';
	import { isAdminPanelLoginPage } from '$elven/tools';
	import { page } from '$app/stores';

	export let isAdmin = false;
</script>

<svelte:head>
	<title>{$page.stuff.title ? $page.stuff.title : 'elven'}</title>
</svelte:head>

<div class="container">
	<Progress />
	{#if isAdmin}
		<Header />
	{/if}

	<main>
		<slot />
	</main>

	<ServiceWrapper bind:isAdmin />
</div>

<style lang="scss">
	.container {
		word-break: break-word;
		min-height: 100vh;
		display: grid;
		grid-template-columns: 1fr;
		grid-template-rows: max-content 1fr max-content;
		main {
			padding: 12px;
		}
	}
</style>
