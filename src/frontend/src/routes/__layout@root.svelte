<script lang="ts">
	import Header from '$lib_oklookat/components/header.svelte';
	import { page } from '$app/stores';
</script>

<svelte:head>
	<meta name="author" content={$page.stuff.author ? $page.stuff.author : 'oklookat'} />
	<title>{$page.stuff.title ? $page.stuff.title : 'oklookat'}</title>
	<meta name="description" content={$page.stuff.description ? $page.stuff.description : 'oklookat home page'} />
</svelte:head>

<div class="container">
	<Header />

	<div class="content">
		<slot />
	</div>
</div>

<style lang="scss">
	.container {
		min-height: 100vh;
		word-break: break-word;
		display: grid;
		grid-template-columns: 1fr;
		// header - content
		grid-template-rows: max-content 1fr;
		gap: 16px;
		padding-bottom: 12px;
		.content {
			height: 100%;
		}
	}
</style>
