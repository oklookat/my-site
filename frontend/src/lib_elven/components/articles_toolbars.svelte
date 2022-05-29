<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Toolbar from '$lib/components/toolbar.svelte';
	import { By, type RAW } from '$elven/types/article';
	import SearchBar from '$elven/components/search_bar.svelte';
	import type { Params, RPH_Event } from '$elven/tools/params';
	import OverlayMobile from '$lib/components/overlay_mobile.svelte';

	export let params: Params<RAW>;

	let searchValue = params.getParam('title') || '';

	if (params.getParam('title')) {
		searchValue = params.getParam('title');
	}

	const dispatch = createEventDispatcher<{
		/** on request param changed */
		paramChanged: RPH_Event<RAW>;
	}>();

	let isFiltersActive = false;

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

{#if isFiltersActive}
	<OverlayMobile title="filters" onClose={() => (isFiltersActive = false)}>
		<div class="filters">
			<div on:click={() => toggleDrafts()}>
				show: {params.getParam('drafts') ? 'drafts' : 'published'}
			</div>

			<div on:click={() => toggleNewest()}>
				age: {params.getParam('newest') ? 'newest' : 'oldest'}
			</div>

			{#if params.getParam('by') === By.updated}
				<div on:click={() => setBy(By.published)}>by: updated date</div>
			{:else if params.getParam('by') === By.published}
				<div on:click={() => setBy(By.created)}>by: published date</div>
			{:else if params.getParam('by') === By.created}
				<div on:click={() => setBy(By.updated)}>by: created date</div>
			{/if}
		</div>
	</OverlayMobile>
{/if}

<div class="toolbars">
	<Toolbar>
		<a href="/elven/articles/create">new</a>
		<div on:click={() => (isFiltersActive = !isFiltersActive)}>filters</div>
	</Toolbar>

	<div class="search">
		<SearchBar bind:value={searchValue} on:search={(e) => search(e.detail)} />
	</div>
</div>

<style lang="scss">
	.filters {
		padding-top: 6px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		gap: 12px;
		div {
			height: 64px;
			padding: 6px;
			background-color: var(--color-level-1);
			cursor: pointer;
			display: flex;
			justify-content: center;
			align-items: center;
		}
	}
	.toolbars {
		height: 104px;
		display: flex;
		flex-direction: column;
		gap: 12px;
		.search {
			height: 50%;
		}
	}
</style>
