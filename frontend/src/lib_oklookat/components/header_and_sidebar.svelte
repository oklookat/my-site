<script lang="ts">
	import { page } from '$app/stores';
	import HamburgerMenu from '$lib/components/hamburger_menu.svelte';
	import Hamburger from '$lib/icons/hamburger.svelte';
	import Navigation from '$oklookat/components/navigation.svelte';

	let isHamburgerOpened = false;
	function toggleHamburger() {
		isHamburgerOpened = !isHamburgerOpened;
	}

	$: onPathChanged($page.url.pathname);
	function onPathChanged(val: string) {
		isHamburgerOpened = false;
	}
</script>

<header class="header">
	<div class="items">
		<div class="hamburger" on:click={() => toggleHamburger()}>
			<div><Hamburger /></div>
		</div>
		<a href="/"><b>oklookat</b></a>
	</div>
	{#if isHamburgerOpened}
		<HamburgerMenu on:closed={() => toggleHamburger()}>
			<Navigation />
		</HamburgerMenu>
	{/if}
</header>

<aside class="sidebar">
	<Navigation />
</aside>

<style lang="scss">
	.header {
		@media screen and(min-width: 755px) {
			display: none;
		}
	}
	.sidebar {
		display: none;
		width: max-content;
		top: 0;
		position: sticky;
		height: 100vh;
		max-height: 100vh;
		background-color: var(--color-level-1);
		@media screen and(min-width: 755px) {
			display: block;
		}
	}
</style>
