<script lang="ts">
	import Utils from '$lib_elven/tools';
	import ElvenLink from '$lib_elven/components/elven_link.svelte';
</script>

<svelte:head>
	<title>{Utils.setTitleElven('home')}</title>
</svelte:head>

<div class="elven-hello base__container">
	<ElvenLink path="/settings">settings</ElvenLink>
</div>
