<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Toolbar from '$lib/components/toolbar.svelte';
	import ToolbarBig from '$lib/components/toolbar_big.svelte';
	import { By, type Article } from '$lib_elven/types/articles';
	import SearchBar from '$lib/components/search_bar.svelte';
	import type { Params, RPH_Event } from '$lib_elven/tools/params';
	import { t } from '$lib/locale';

	export let params: Params<Article>;

	let searchValue = params.getParam('title') || '';

	if (params.getParam('title')) {
		searchValue = params.getParam('title');
	}

	const dispatch = createEventDispatcher<{
		/** on request param changed */
		paramChanged: RPH_Event<Article>;
	}>();

	/** set 'by' param and get articles */
	function setBy(by: By = By.published) {
		dispatch('paramChanged', { name: 'by', val: by });
	}

	/** set 'newest' param and get articles */
	function toggleNewest() {
		dispatch('paramChanged', {
			name: 'newest',
			val: !params.getParam('newest')
		});
	}

	/** set 'drafts' param and get articles */
	function toggleDrafts() {
		dispatch('paramChanged', {
			name: 'drafts',
			val: !params.getParam('drafts')
		});
	}

	/** search by title */
	function search(val: string) {
		dispatch('paramChanged', { name: 'title', val: val });
	}
</script>

<div class="toolbars">
	<ToolbarBig>
		<a href="/elven/articles/create">{$t('elven.articles.create')}</a>
	</ToolbarBig>

	<div class="search">
		<SearchBar bind:value={searchValue} on:search={(e) => search(e.detail)} />
	</div>

	<Toolbar>
		<div>
			<div class="pointer" on:click={() => toggleDrafts()}>
				{params.getParam('drafts')
					? $t('elven.articles.drafts')
					: $t('elven.articles.published')}
			</div>
		</div>

		<div>
			<div class="pointer" on:click={() => toggleNewest()}>
				{params.getParam('newest')
					? $t('elven.articles.newest')
					: $t('elven.articles.oldest')}
			</div>
		</div>

		<div>
			{#if params.getParam('by') === By.updated}
				<div class="pointer" on:click={() => setBy(By.published)}>
					{$t('elven.articles.byUpdatedDate')}
				</div>
			{:else if params.getParam('by') === By.published}
				<div class="pointer" on:click={() => setBy(By.created)}>
					{$t('elven.articles.byPublishedDate')}
				</div>
			{:else if params.getParam('by') === By.created}
				<div class="pointer" on:click={() => setBy(By.updated)}>
					{$t('elven.articles.byCreatedDate')}
				</div>
			{/if}
		</div>
	</Toolbar>
</div>

<style lang="scss">
	.toolbars {
		display: flex;
		flex-direction: column;
		gap: 12px;
		width: 100%;
		.search {
			height: 54px;
		}
	}
</style>
