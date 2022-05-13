<script lang="ts">
	import CFile from '$lib/components/elven/file.svelte';
	import type { Items } from '$lib/types';
	import type { File } from '$lib/types/files';
	import { createEventDispatcher } from 'svelte';

	export let items: Items<File>;

	const dispatch = createEventDispatcher<{
		deleted: number;
	}>();

	async function whenDeleted(counter: number) {
		delete items.data[counter];
		items = items;
		dispatch('deleted', counter);
	}
</script>

<div class="list">
	{#if items && items.data}
		{#each Object.entries(items.data) as [counter, file]}
			<CFile bind:file onDeleted={() => whenDeleted(parseInt(counter, 10))} />
		{/each}
	{/if}
</div>
