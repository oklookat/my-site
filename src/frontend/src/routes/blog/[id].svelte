<script lang="ts">
	import { browser } from '$app/env';
	import { getParser } from '$lib/tools/markdown';

	import ArticleCover from '$lib_elven/components/article_cover.svelte';
	import type { Article } from '$lib_elven/types/articles';

	export let article: Article;

	let parseMarkdown: (data: string) => string;

	if (browser) {
		parseMarkdown = getParser();
	}
</script>

<svelte:head>
	<title>{article.title}</title>
</svelte:head>

<article>
	<div class="article__content">
		<h1>{article.title}</h1>
		<ArticleCover {article} />
		{#if parseMarkdown}
			{@html parseMarkdown(article.content)}
		{/if}
	</div>
</article>

<style lang="scss">
	article {
		font-size: 1.3rem;
		width: 95%;
		max-width: 624px;
		margin: auto;
	}
</style>
