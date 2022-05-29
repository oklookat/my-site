<script lang="ts">
	import Header from '$oklookat/components/header.svelte';
	import { page } from '$app/stores';
</script>

<svelte:head>
	<meta name="author" content={$page.stuff.author ? $page.stuff.author : 'oklookat'} />
	<title>{$page.stuff.title ? $page.stuff.title : 'oklookat'}</title>
	<meta
		name="description"
		content={$page.stuff.description ? $page.stuff.description : 'oklookat home page'}
	/>
</svelte:head>

<div class="container">
	<Header />

	<main>
		<slot />
	</main>
</div>

<style lang="scss">
	.container {
		word-break: break-word;
		min-height: 100vh;
		display: grid;
		grid-template-columns: 1fr;
		grid-template-rows: max-content 1fr;
		main {
			padding: 12px;
		}
	}
</style>
