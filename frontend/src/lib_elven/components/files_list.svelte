<script lang="ts">
	import CFile from '$elven/components/file.svelte';
	import type { Items } from '$elven/types';
	import type { File } from '$elven/types/file';
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

{#if items && items.data}
	{#each Object.entries(items.data) as [counter, file]}
		<CFile bind:file onDeleted={() => whenDeleted(parseInt(counter, 10))} />
	{/each}
{/if}
