<script lang="ts">
	import { page } from '$app/stores';
	import { onDestroy } from 'svelte';
	

	let isUnknown = false;
	let isArticles = false;
	let isFiles = false;

	const unsub = page.subscribe((v) => {
		if (!v || !v.url || !v.url.pathname) {
			return;
		}
		onPathChanged(v.url.pathname);
	});

	onDestroy(() => {
		unsub();
	});

	function onPathChanged(path: string) {
		if (!path) {
			isUnknown = true;
			return;
		}
		isFiles = path.includes('/elven/files');
		isArticles = path.includes('/elven/articles');
		isUnknown = !isFiles && !isArticles;
	}
</script>

<nav class="header base__container">
	<div class="header__items">
		<a href="/elven">
			<div class={isUnknown ? 'base__links' : ''}>
				<div>elven</div>
			</div></a
		>

		<a href="/elven/articles">
			<div class={isArticles ? 'base__links' : ''}>
				<div>articles</div>
			</div>
		</a>

		<a href="/elven/files">
			<div class={isFiles ? 'base__links' : ''}>
				<div>files</div>
			</div>
		</a>
	</div>
</nav>

<style lang="scss">
	.header {
		font-weight: bold;
		height: max-content;
		display: flex;
		&__items {
			padding-top: 14px;
			padding-bottom: 14px;
			width: 100%;
			display: flex;
			flex-wrap: wrap;
			justify-content: center;
			gap: 42px;
		}
	}
</style>
