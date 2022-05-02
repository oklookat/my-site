<script lang="ts">
	import SearchBar from '$lib/components/search_bar.svelte';
	import Toolbar from '$lib/components/toolbar.svelte';
	import type { Params, RPH_Event } from '$lib_elven/tools/params';
	import type { Article } from '$lib_elven/types/articles';
	import { createEventDispatcher } from 'svelte';
	import { _ } from 'svelte-i18n';

	export let params: Params<Article>;

	let searchValue = '';

	if (params.getParam('title')) {
		searchValue = params.getParam('title');
	}

	const dispatch = createEventDispatcher<{
		/** on request param changed */
		paramChanged: RPH_Event<Article>;
	}>();

	/** set 'newest' param and get articles */
	function toggleNewest() {
		dispatch('paramChanged', { name: 'newest', val: !params.getParam('newest') });
	}

	/** search by title */
	function search(val: string) {
		dispatch('paramChanged', { name: 'title', val: val });
	}
</script>

<div class="toolbars">
	<div class="search">
		<SearchBar bind:value={searchValue} on:search={(e) => search(e.detail)} />
	</div>

	<Toolbar>
		<div>
			<div class="pointer" on:click={() => toggleNewest()}>
				{params.getParam('newest')
					? $_('elven.components.articlesToolbars.newest')
					: $_('elven.components.articlesToolbars.oldest')}
			</div>
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
