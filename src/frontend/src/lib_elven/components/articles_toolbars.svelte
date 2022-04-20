<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	// ui
	import Toolbar from '$lib_elven/components/toolbar.svelte';
	import ToolbarBig from '$lib_elven/components/toolbar_big.svelte';
	import ElvenLink from '$lib_elven/components/elven_link.svelte';
	// article
	import { By } from '$lib_elven/types/articles';
	import type { Params } from '$lib_elven/types/articles';
	import SearchBar from '$lib_elven/components/search_bar.svelte';

	export let params: Params;

	let searchValue = '';

	$: onTitleParamChanged(params.title);
	function onTitleParamChanged(value: string | undefined) {
		searchValue = value;
	}

	const dispatch = createEventDispatcher<{
		/** on request param changed */
		paramChanged: { name: string; val: string | boolean };
	}>();

	/** set 'by' param and get articles */
	function setBy(by: By = By.published) {
		dispatch('paramChanged', { name: 'by', val: by });
	}

	/** set 'start' param and get articles */
	function toggleNewest() {
		dispatch('paramChanged', { name: 'newest', val: !params.newest });
	}

	/** set 'show' param and get articles */
	function togglePublished() {
		if (typeof params.published !== 'boolean') {
			console.error(
				'params.published not a boolean, possibly toolbar strange behavior. Recheck your type conversion'
			);
		}
		dispatch('paramChanged', { name: 'published', val: !params.published });
	}

	/** search by filename */
	function search(val: string) {
		params.title = val;
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
			<div class="pointer" on:click={() => togglePublished()}>
				{params.published ? 'published' : 'drafts'}
			</div>
		</div>
		<div>
			<div class="pointer" on:click={() => toggleNewest()}>
				{params.newest ? 'newest' : 'oldest'}
			</div>
		</div>
		<div>
			{#if params.by === By.updated}
				<div class="pointer" on:click={() => setBy(By.published)}>by updated date</div>
			{:else if params.by === By.published}
				<div class="pointer" on:click={() => setBy(By.created)}>by published date</div>
			{:else if params.by === By.created}
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
