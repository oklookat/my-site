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
	import Progress from '$elven/plugins/elvenProgress/progress.svelte';
	import { isAdminPanelLoginPage } from '$elven/tools';
	import { page } from '$app/stores';
	import HeaderAndSidebar from '$elven/components/header_and_sidebar.svelte';
	import Notify from '$elven/plugins/elvenNotify/notify.svelte';
	import Player from '$elven/plugins/elvenPlayer/player.svelte';
	import Confirm from '$elven/plugins/elvenChoose/confirm.svelte';

	export let isAdmin = false;
</script>

<svelte:head>
	<title>{$page.stuff.title ? $page.stuff.title : 'elven'}</title>
</svelte:head>

<div class="container">
	<Progress />
	{#if isAdmin}
		<div class="header">
			<HeaderAndSidebar />
		</div>
	{/if}

	<main>
		<slot />
	</main>

	<div class="service">
		<Notify />
		{#if isAdmin}
			<Confirm />
			<Player />
		{/if}
	</div>
</div>

<style lang="scss">
	.container {
		word-break: break-word;
		min-height: 100vh;
		display: grid;
		// when header
		grid-template-columns: 1fr 1fr;
		grid-template-rows: max-content 1fr auto;
		grid-template-areas:
			'header header'
			'main main'
			'service service';
		@media screen and(min-width: 755px) {
			// when sidebar
			grid-template-columns: max-content 1fr;
			grid-template-rows: 1fr 1fr auto;
			grid-template-areas:
				'header main'
				'header main'
				'service service';
		}

		.header {
			grid-area: header;
		}
		main {
			grid-area: main;
			padding: 12px;
		}
		.service {
			grid-area: service;
			display: flex;
			height: fit-content;
			width: 100%;
			z-index: 100;
			position: sticky;
			bottom: 0;
		}
	}
</style>
