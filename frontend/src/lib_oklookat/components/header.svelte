<script lang="ts">
	import { page, session } from '$app/stores';
	import HamburgerMenu from '$lib/components/hamburger_menu.svelte';
	import Hamburger from '$lib/icons/hamburger.svelte';
	import Link from '$oklookat/components/link.svelte';

	let isHamburgerOpened = false;
	function toggleHamburger() {
		isHamburgerOpened = !isHamburgerOpened;
	}

	$: onPathChanged($page.url.pathname);
	function onPathChanged(val: string) {
		isHamburgerOpened = false;
	}
</script>

{#if isHamburgerOpened}
	<HamburgerMenu on:closed={() => toggleHamburger()}>
		<div class="navigateme">
			<nav>
				<Link path="/">дом</Link>
				<Link path="/blog">блог</Link>
				<Link path="/toys">штуки</Link>
			</nav>
			{#if $session.user.isAdmin}
				<hr />
				<nav>
					<Link path="/elven">elven</Link>
				</nav>
			{/if}
		</div>
	</HamburgerMenu>
{/if}

<header>
	<div class="items">
		<div class="hamburger" on:click={() => toggleHamburger()}>
			<div><Hamburger /></div>
		</div>
		<a href="/"><b>oklookat</b></a>
	</div>
</header>
