<script lang="ts">
	import { dateToReadable } from '$lib/tools/dates';
	import type { RAW } from '$lib/types/article';
	import { createEventDispatcher } from 'svelte';
	import { t } from '$lib/locale';
import Actions from '$lib/components/elven/article_actions.svelte';
import Cover from '$lib/components/elven/article_cover.svelte';

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
	<Actions
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
				<Cover {article} />
			</div>
		{/if}
		<div class="title">
			{article.title}
		</div>
	</div>
</article>
