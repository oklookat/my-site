<script context="module" lang="ts">
	export const load: Load = async (e) => {
		let resp: Response | undefined = undefined;
		const networkArticle = new NetworkArticle('', e.fetch);
		let statusCode = 200;

		const stuff = e.stuff;
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
			const art = (await resp.json()) as RAW;
			output.props.article = art;
			stuff.title = art.title;
			stuff.description = cutString(art.content, 94);
		}

		return output;
	};
</script>

<script lang="ts">
	import { browser } from '$app/env';
	import { cutString } from '$elven/tools';
	import { getParser } from '$elven/tools/markdown';

	import ArticleCover from '$elven/components/article_cover.svelte';
	import NetworkArticle from '$elven/network/article';
	import type { RAW } from '$elven/types/article';
	import type { Load, LoadOutput } from '@sveltejs/kit';

	export let article: RAW;

	let parseMarkdown: (data: string) => string;

	if (browser) {
		parseMarkdown = getParser();
	}
</script>

{#if article}
	<article>
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
