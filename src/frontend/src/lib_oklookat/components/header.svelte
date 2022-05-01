<script lang="ts">
	import { goto } from '$app/navigation';

	import { session } from '$app/stores';
	import HamburgerMenu from '$lib/components/hamburger_menu.svelte';
	import Hamburger from '$lib/icons/hamburger.svelte';

	let isHamburgerOpened = false;
	function toggleHamburger() {
		isHamburgerOpened = !isHamburgerOpened;
	}

	async function goMain() {
		isHamburgerOpened = false;
		await goto('/');
	}
</script>

<div class="header base__container">
	<div class="items">
		<div class="hamburger" on:click={() => toggleHamburger()}>
			<Hamburger />
		</div>
		<b class="logo" on:click={() => goMain()}>oklookat</b>
	</div>
</div>

{#if isHamburgerOpened}
	<HamburgerMenu on:closed={() => toggleHamburger()}>
		{#if $session.user.isAdmin}
			<a href="/elven" on:click={() => toggleHamburger()}>elven</a>
		{/if}
		<a href="/" on:click={() => toggleHamburger()}>дом</a>
		<a href="/blog" on:click={() => toggleHamburger()}>блог</a>
		<a href="/toys" on:click={() => toggleHamburger()}>штуки</a>
	</HamburgerMenu>
{/if}

<style lang="scss">
	.header {
		height: 44px;
		align-items: center;
		.items {
			height: 100%;
			width: 100%;

			* {
				height: 100%;
				display: flex;
				align-items: center;
				cursor: pointer;
			}

			display: grid;
			grid-template-rows: 1fr;
			grid-template-columns: max-content 1fr;
			gap: 14px;

			.hamburger {
				width: 24px;
			}
			.logo {
				width: 100%;
			}
		}
	}
</style>
