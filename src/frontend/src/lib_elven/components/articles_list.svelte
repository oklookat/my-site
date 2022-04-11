<script lang="ts">
	import type { Data } from '$lib_elven/types';
	import type { Article } from '$lib_elven/types/articles';
	import CArticle from './article.svelte';

	export let items: Data<Article>;
	export let onDeleted: () => void;

	/** delete file from files array */
	function deleteFromArray(counter: number) {
		delete items.data[counter];
		items = items;
	}

	function whenDeleted(counter: number) {
		deleteFromArray(counter);
		onDeleted();
	}
</script>

<div class="list">
	{#if items && items.data}
		{#each Object.entries(items.data) as [counter, article]}
			<CArticle {article} onDeleted={() => whenDeleted(parseInt(counter, 10))} />
		{/each}
	{/if}
</div>
