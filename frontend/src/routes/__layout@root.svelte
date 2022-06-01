<script lang="ts">
	import { page } from '$app/stores';
	import HeaderAndSidebar from '$oklookat/components/header_and_sidebar.svelte';
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
	<HeaderAndSidebar />

	<main>
		<slot />
	</main>
</div>

<style lang="scss">
	.container {
		word-break: break-word;
		min-height: 100vh;
		display: grid;
		// when header
		grid-template-columns: 1fr;
		grid-template-rows: max-content 1fr;
		@media screen and(min-width: 755px) {
			// when sidebar
			grid-template-columns: max-content 1fr;
			grid-template-rows: 1fr;
		}
		main {
			padding: 12px;
		}
	}
</style>
