<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';

	export const load: Load = async (event) => {
		const isError = event.session.user.isError;
		const isAdmin = event.session.user.isAdmin;
		const isLoginPage = isAdminPanelLoginPage(event.url);

		if (isError) {
			// redirect to main page if auth check error
			return {
				status: 302,
				redirect: ''
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
			props: {
				isAdmin: isAdmin
			}
		};
	};
</script>

<script lang="ts">
	// components
	import Header from '$lib_elven/components/header.svelte';
	import ServiceWrapper from '$lib_elven/components/service_wrapper.svelte';
	import Progress from '$lib_elven/plugins/elvenProgress/progress.svelte';
	import { isAdminPanelLoginPage } from '$lib/tools';

	export let isAdmin = false;
</script>

<div class="container">
	<Progress />

	{#if isAdmin}
		<Header />
	{/if}

	<div class="content">
		<slot />
	</div>

	<ServiceWrapper bind:isAdmin />
</div>

<style lang="scss">
	.container {
		min-height: 100vh;
		word-break: break-word;
		display: grid;
		grid-template-columns: 1fr;
		// header - content - service
		grid-template-rows: max-content 1fr min-content;
		gap: 16px;

		.content {
			height: 100%;
			width: 100%;
			font-size: 1.1rem;
			line-height: 1.46rem;
			letter-spacing: 0.0007rem;
		}
	}
</style>
