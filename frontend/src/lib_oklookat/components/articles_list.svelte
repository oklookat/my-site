<script lang="ts">
	import type { Items } from '$elven/types';
	import type { RAW } from '$elven/types/article';
	import CArticle from './article.svelte';

	export let items: Items<RAW>;
</script>

{#if items && items.data}
	{#each Object.entries(items.data) as [counter, article]}
		<CArticle {article} />
	{/each}
{/if}
