<script lang="ts">
	import Header from '$lib_oklookat/components/header.svelte';
</script>

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
