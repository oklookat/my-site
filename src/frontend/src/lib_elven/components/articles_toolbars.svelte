<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Toolbar from '$lib_elven/components/toolbar.svelte';
	import ToolbarBig from '$lib_elven/components/toolbar_big.svelte';
	import ElvenLink from '$lib_elven/components/link.svelte';
	import { By, type Article } from '$lib_elven/types/articles';
	import SearchBar from '$lib_elven/components/search_bar.svelte';
	import type { Params, RPH_Event } from '$lib_elven/tools/params';

	export let params: Params<Article>;

	let searchValue = '';

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
		<ElvenLink path="/articles/create">new</ElvenLink>
	</ToolbarBig>

	<div class="search">
		<SearchBar bind:value={searchValue} on:search={(e) => search(e.detail)} placeholder="search" />
	</div>

	<Toolbar>
		<div>
			<div class="pointer" on:click={() => toggleDrafts()}>
				{params.getParam('drafts') ? 'drafts' : 'published'}
			</div>
		</div>
		<div>
			<div class="pointer" on:click={() => toggleNewest()}>
				{params.getParam('newest') ? 'newest' : 'oldest'}
			</div>
		</div>
		<div>
			{#if params.getParam('by') === By.updated}
				<div class="pointer" on:click={() => setBy(By.published)}>by updated date</div>
			{:else if params.getParam('by') === By.published}
				<div class="pointer" on:click={() => setBy(By.created)}>by published date</div>
			{:else if params.getParam('by') === By.created}
				<div class="pointer" on:click={() => setBy(By.updated)}>by created date</div>
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
