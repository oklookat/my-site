<script lang="ts">
	import { dateToReadable } from '$elven/tools/dates';
	import type { RAW } from '$elven/types/article';
	import ArticleCover from './article_cover.svelte';
	import ArticleActions from '$elven/components/article_actions.svelte';
	import { createEventDispatcher } from 'svelte';
	import { t } from '$lib/locale';

	export let article: RAW;

	const dispatch = createEventDispatcher<{
		/** on article deleted */
		deleted: void;
	}>();

	function convertDate(date: string | number | Date): string {
		return dateToReadable(date);
	}

	/** is article selected? (actions menu/overlay opened) */
	let isSelected = false;
	let selectedMouseEvent: MouseEvent;
	function onSelected(e: MouseEvent) {
		selectedMouseEvent = e;
		isSelected = true;
	}

	function onDeleted() {
		dispatch('deleted');
		isSelected = false;
	}
</script>

{#if isSelected}
	<ArticleActions
		{article}
		mouseEvent={selectedMouseEvent}
		onDisabled={() => (isSelected = false)}
		onDeleted={() => {
			onDeleted();
		}}
	/>
{/if}

<article class="article" on:click={(e) => onSelected(e)}>
	<div class="meta">
		<div class="meta__item">
			{#if article.is_published && article.published_at}
				{convertDate(article.published_at)}
			{:else if article.updated_at}
				{convertDate(article.updated_at)}
			{/if}
		</div>

		<div class="meta__item">
			{#if article.is_published}
				{$t('elven.articles.published')}
			{:else}
				{$t('elven.articles.draft')}
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
