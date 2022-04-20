<script lang="ts">
	// utils
	import Dates from '$lib_elven/tools/dates';
	// article
	import type { Article } from '$lib_elven/types/articles';
	import ArticleCover from './article_cover.svelte';
	import ArticleActions from '$lib_elven/components/article_actions.svelte';

	export let article: Article;

	/** on article deleted */
	export let onDeleted: () => void;

	function convertDate(date: string | number | Date): string {
		return Dates.convert(date);
	}

	/** is article selected? (actions menu/overlay opened) */
	let isSelected = false;
	let selectedMouseEvent: MouseEvent;
	function onSelected(e: MouseEvent) {
		selectedMouseEvent = e;
		isSelected = true;
	}
</script>

{#if isSelected}
	<ArticleActions
		{article}
		mouseEvent={selectedMouseEvent}
		onDisabled={() => (isSelected = false)}
		onDeleted={() => onDeleted()}
	/>
{/if}

<article class="article" on:click={(e) => onSelected(e)}>
	<div class="meta">
		<div class="meta__item">
			{#if article.is_published && article.published_at}
				{convertDate(article.published_at)}
			{:else}
				{convertDate(article.updated_at)}
			{/if}
		</div>

		<div class="meta__item">
			{#if article.is_published}
				published
			{:else}
				draft
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

<style lang="scss">
	.article {
		border-radius: var(--border-radius);
		background-color: var(--color-level-1);
		cursor: pointer;
		width: 100%;
		display: flex;
		flex-direction: column;

		.meta {
			padding: 12px 12px 8px;
			display: flex;
			flex-direction: row;
			height: fit-content;
			gap: 12px;
			flex-wrap: wrap;

			.meta__item {
				font-size: 0.9rem;
				min-width: 44px;
				border-radius: 12px;
				min-height: 32px;
				height: fit-content;
				display: flex;
				justify-content: center;
				align-items: center;
			}
		}

		.main {
			width: 100%;
			height: 100%;
			display: flex;
			flex-direction: column;
			gap: 4px;

			.title {
				font-size: 1.3rem;
				line-height: 2rem;
				letter-spacing: 0.0099rem;
				padding: 8px 12px 18px;
			}
		}
	}
</style>
