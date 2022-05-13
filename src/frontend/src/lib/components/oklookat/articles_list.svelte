<script lang="ts">
	import type { Items } from '$lib/types';
	import type { Article } from '$lib/types/articles';
	import CArticle from './article.svelte';

	export let items: Items<Article>;
</script>

<div class="list">
	{#if items && items.data}
		{#each Object.entries(items.data) as [counter, article]}
			<CArticle {article}  />
		{/each}
	{/if}
</div>