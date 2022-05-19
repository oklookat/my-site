<script context="module" lang="ts">
	export const load: Load = async (e) => {
		let resp: Response | undefined = undefined;
		const networkArticle = new NetworkArticle('', e.fetch);
		let statusCode = 200;

		const stuff = e.stuff
		const output: LoadOutput = {
			status: statusCode,
			stuff: stuff,
			props: {
				article: null
			}
		};

		try {
			resp = await networkArticle.get(e.params.id);
			statusCode = resp.status;
		} catch (err) {}

		if (resp && resp.ok && output.props) {
			const art = (await resp.json()) as Article;
			output.props.article = art;
			stuff.title = art.title;
			stuff.description = cutString(art.content, 94);
		}

		return output;
	};
</script>

<script lang="ts">
	import { browser } from '$app/env';
	import { cutString } from '$lib/tools';
	import { getParser } from '$lib/tools/markdown';

	import ArticleCover from '$lib/components/elven/article_cover.svelte';
	import NetworkArticle from '$lib/network/article';
	import type { Article } from '$lib/types/articles';
	import type { Load, LoadOutput } from '@sveltejs/kit';

	export let article: Article;

	let parseMarkdown: (data: string) => string;

	if (browser) {
		parseMarkdown = getParser();
	}
</script>

{#if article}
	<article class="base__container">
		<div class="article__content">
			<h1>{article.title}</h1>
			<ArticleCover {article} />
			{#if parseMarkdown}
				{@html parseMarkdown(article.content)}
			{/if}
		</div>
	</article>
{/if}

<style lang="scss">
	article {
		font-size: 1.3rem;
	}
</style>
