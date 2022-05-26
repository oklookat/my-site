<script lang="ts">
	import { page } from '$app/stores';
	import Link from '$elven/components/link.svelte';
	import HamburgerMenu from '$lib/components/hamburger_menu.svelte';
	import Hamburger from '$lib/icons/hamburger.svelte';

	let isHamburgerOpened = false;
	function toggleHamburger() {
		isHamburgerOpened = !isHamburgerOpened;
	}

	$: onPathChanged($page.url.pathname);
	function onPathChanged(val: string) {
		isHamburgerOpened = false;
	}
</script>

<header>
	<div class="hamburger" on:click={() => toggleHamburger()}>
		<Hamburger />
	</div>
</header>

{#if isHamburgerOpened}
	<HamburgerMenu on:closed={() => (isHamburgerOpened = false)}>
		<div class="navigateme">
			<nav>
				<Link path="/elven">elven</Link>
				<Link path="/elven/articles">articles</Link>
				<Link path="/elven/files">files</Link>
			</nav>

			<hr />

			<nav>
				<Link path="/">main site</Link>
				<Link path="/elven/settings">settings</Link>
			</nav>
		</div>
	</HamburgerMenu>
{/if}

<style lang="scss">
	header {
		position: sticky;
		top: 0;
		z-index: 99;
		height: 44px;
		width: 100%;
		//
		padding: 6px;
		background-color: var(--color-level-1);
		display: flex;
		align-items: center;
		* {
			height: 100%;
			display: flex;
			align-items: center;
			cursor: pointer;
		}

		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
		gap: 14px;

		.hamburger {
			height: 100%;
			width: 24px;
		}
	}

	.navigateme {
		background-color: var(--color-level-1);
		font-weight: bold;
		height: 100%;
		min-width: 94px;
		width: max-content;
		&,
		nav {
			display: flex;
			flex-wrap: wrap;
			flex-direction: column;
		}
		hr,
		nav {
			width: 100%;
		}
		nav {
			// nav item
			:global(*) {
				min-height: 44px;
				height: max-content;
				width: 100%;
				display: flex;
				justify-content: center;
				align-items: center;
				&:hover {
					background-color: var(--color-level-2);
				}
			}
		}
	}
</style>
