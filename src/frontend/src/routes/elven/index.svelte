<script lang="ts">
	import ElvenLink from '$lib_elven/components/elven_link.svelte';
	import { setTitleElven } from '$lib_elven/tools';
</script>

<svelte:head>
	<title>{setTitleElven('home')}</title>
</svelte:head>

<div class="elven-hello base__container">
	<ElvenLink path="/settings">settings</ElvenLink>
</div>
