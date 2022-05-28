<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';

	export const load: Load = async (event) => {
		return {
			props: {
				statusCode: event.status
			}
		};
	};
</script>

<script lang="ts">
	import Error400 from '$oklookat/components/error400.svelte';
	import Error404 from '$oklookat/components/error404.svelte';
	import Error500 from '$oklookat/components/error500.svelte';
	import ErrorInternet from '$oklookat/components/errorInternet.svelte';
	import ErrorUnknown from '$oklookat/components/errorUnknown.svelte';
	import { onMount } from 'svelte';

	export let statusCode: number | null;

	let isOffline = false;
	onMount(() => {
		isOffline = !navigator.onLine;
		if (isOffline) {
			statusCode = null;
		}
	});
</script>

<div class="ohno">
	<div class="message">
		{#if isOffline}
			<ErrorInternet />
		{:else if statusCode}
			{#if statusCode === 404}
				<Error404 />
			{:else if statusCode === 500}
				<Error500 />
			{:else if statusCode > 399 && statusCode < 500}
				<Error400 />
			{:else}
				<ErrorUnknown />
			{/if}
		{/if}
	</div>
	<div class="links">
		<a class="gohome" href="/">на главную</a>
		<a class="tg" href="tg://resolve?domain=andget">написать в telegram</a>
	</div>
</div>

<style lang="scss">
	.ohno {
		min-height: 100%;
		margin: 12px;
		display: grid;
		grid-template-rows: 1fr auto;
		align-items: center;
		justify-content: center;
		gap: 24px;
		.links {
			align-self: flex-end;
			display: flex;
			flex-direction: row;
			align-items: center;
			justify-content: center;
			gap: 12px;
			a {
				font-weight: bold;
			}
			.gohome,
			.tg {
				border-bottom: 2px solid red;
			}
		}
	}
</style>
