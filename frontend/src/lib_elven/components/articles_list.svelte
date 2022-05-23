<script lang="ts">
	import type { Items } from '$elven/types';
	import type { RAW } from '$elven/types/article';
	import { createEventDispatcher } from 'svelte';
	import CArticle from './article.svelte';

	export let items: Items<RAW>;

	const dispatch = createEventDispatcher<{
		/** on article deleted (with counter) */
		deleted: number;
	}>();

	/** delete file from files array */
	function deleteFromArray(counter: number) {
		delete items.data[counter];
		items = items;
	}

	function onDeleted(counter: number) {
		deleteFromArray(counter);
		dispatch('deleted', counter);
	}
</script>

<div class="list">
	{#if items && items.data}
		{#each Object.entries(items.data) as [counter, article]}
			<CArticle {article} on:deleted={() => onDeleted(parseInt(counter, 10))} />
		{/each}
	{/if}
</div>
