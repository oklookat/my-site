<script context="module">
	export const load = async (event) => {
		if (Validator.isAdminPanelLoginPage(event.url)) {
			return {};
		}
		const isAdmin = !!event.session.user && !!event.session.user.isAdmin;
		if (!isAdmin) {
			return {
				status: 302,
				redirect: '/elven/login'
			};
		}
		return {};
	};
</script>

<script lang="ts">
	import { browser } from '$app/env';
	import { page } from '$app/stores';
	// main style
	import '$lib_elven/assets/global.scss';
	// plugins
	import Progress from '$lib_elven/plugins/elvenProgress/progress.svelte';
	// components
	import Header from '$lib_elven/components/header.svelte';
	import ServiceWrapper from '$lib_elven/components/service_wrapper.svelte';
	import Validator from '$lib_elven/validators';
	import Utils from '$lib_elven/tools';

	// TODO: handle auth on endpoint. Like send auth request in endpoint, and if 403, not send html.
	// state
	// let isAuthPage = false;
	// location.subscribe((value) => {
	//   isAuthPage = value.includes("/login") || value.includes("/logout");
	// });
	// let is404Page = false;
	// GlobalState.isNotFoundPage.subscribe((value) => {
	//   is404Page = value;
	// });

	let isNotAuthPage = false;
	page.subscribe((val) => {
		const pathname = val.url.pathname;
		isNotAuthPage = pathname !== '/elven/login' && pathname !== '/elven/logout';
	});
</script>

<div class="container">
	{#if isNotAuthPage && browser}
		<Progress />
		<Header />
	{/if}
	<div class="content">
		<slot />
	</div>
	{#if isNotAuthPage && browser}
		<ServiceWrapper />
	{/if}
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
	}

	.content {
		height: 100%;
		width: 100%;
		font-size: 1.1rem;
		line-height: 1.46rem;
		letter-spacing: 0.0007rem;
	}
</style>
