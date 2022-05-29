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
	<div class="items">
		<div class="hamburger" on:click={() => toggleHamburger()}>
			<div><Hamburger /></div>
		</div>
		<a href="/elven"><b>elven</b></a>
	</div>
</header>

{#if isHamburgerOpened}
	<HamburgerMenu on:closed={() => (isHamburgerOpened = false)}>
		<div class="navigateme">
			<nav>
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
