<script lang="ts">
	import { session } from '$app/stores';
	import Link from '$oklookat/components/link.svelte';
</script>

<div class="navigateme">
	<nav>
		<Link path="/">дом</Link>
		<Link path="/blog">блог</Link>
		<Link path="/toys">штуки</Link>
		{#if $session.user.isAdmin}
			<hr />
			<Link path="/elven">elven</Link>
		{/if}
	</nav>
</div>
