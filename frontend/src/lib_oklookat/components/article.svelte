<script lang="ts">
	import { goto } from '$app/navigation';
	import ArticleCover from '$elven/components/article_cover.svelte';
	import { dateToReadable } from '$elven/tools/dates';

	import type { RAW } from '$elven/types/article';

	export let article: RAW;

	function convertDate(date: string | number | Date): string {
		return dateToReadable(date);
	}

	async function goToArticle() {
		await goto(`/blog/${article.id}`);
	}
</script>

<article class="article" on:click={async () => await goToArticle()}>
	<div class="meta">
		<div class="meta__item">
			{#if article.is_published && article.published_at}
				{convertDate(article.published_at)}
			{/if}
		</div>
	</div>
	<div class="main">
		{#if article.cover_id && article.cover_path && article.cover_extension}
			<div class="cover">
				<ArticleCover {article} />
			</div>
		{/if}
		<div class="title">
			{article.title}
		</div>
	</div>
</article>
