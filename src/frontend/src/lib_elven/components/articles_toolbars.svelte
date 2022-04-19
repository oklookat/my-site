<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	// ui
	import Toolbar from '$lib_elven/components/toolbar.svelte';
	import ToolbarBig from '$lib_elven/components/toolbar_big.svelte';
	import ElvenLink from '$lib_elven/components/elven_link.svelte';
	// article
	import { By } from '$lib_elven/types/articles';
	import type { Params } from '$lib_elven/types/articles';
	import type { Category } from '$lib_elven/types/articles/categories';
	import CategoriesSelector from '$lib_elven/components/categories_selector.svelte';
	import type { Counter } from '$lib_elven/types';
	import { ToolsCategories } from '$lib_elven/tools/categories';
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

	const ALL_CATEGORIES = '-1';
	const WITHOUT_CATEGORY = '-2';

	/** custom categories for searching by category */
	const customCategories: Record<number, Category> = {};
	customCategories[ALL_CATEGORIES] = {
		id: ALL_CATEGORIES,
		name: 'All'
	};
	customCategories[WITHOUT_CATEGORY] = {
		id: WITHOUT_CATEGORY,
		name: 'Without category'
	};

	/** default selected */
	let selected = ALL_CATEGORIES;

	/** is ready to display toolbars? */
	let isLoaded = false;

	onMount(async () => {
		if (!params.without_category) {
			selected = ALL_CATEGORIES;
		} else {
			selected = WITHOUT_CATEGORY;
			isLoaded = true;
			return;
		}

		if (!params.category_name) {
			isLoaded = true;
			return;
		}

		// if we have category name in params, we need to get id of this category
		try {
			const counter = await ToolsCategories.getCounterByName('', params.category_name);
			selected = counter;
		} catch (err) {}

		isLoaded = true;
	});

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

	/** sort by category */
	function onCategoryChanged(selected: { counter: Counter; category: Category }) {
		if (selected.counter === ALL_CATEGORIES || selected.counter === WITHOUT_CATEGORY) {
			// if we search without category
			// category_name a priori should be false/undefined
			params.category_name = undefined;
			dispatch('paramChanged', { name: 'category_name', val: params.category_name });

			// set
			params.without_category = selected.counter === WITHOUT_CATEGORY;
			dispatch('paramChanged', { name: 'without_category', val: params.without_category });
			return;
		}

		// if we search by category name
		// without_category a priori should be false/undefined
		params.without_category = undefined;
		dispatch('paramChanged', { name: 'without_category', val: params.without_category });

		// set
		params.category_name = selected.category.name;
		dispatch('paramChanged', { name: 'category_name', val: params.category_name });
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
		<ElvenLink path="/articles/categories">categories</ElvenLink>
	</ToolbarBig>

	{#if isLoaded}
		<div class="oneline">
			<div class="one">
				<Toolbar>
					<CategoriesSelector
						{customCategories}
						{selected}
						on:selected={(e) => onCategoryChanged(e.detail)}
					/>
				</Toolbar>
			</div>
			<div class="two">
				<SearchBar
					bind:value={searchValue}
					on:search={(e) => search(e.detail)}
					placeholder="search"
				/>
			</div>
		</div>
	{/if}

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
		.oneline {
			display: grid;
			grid-template-rows: 1fr;
			grid-template-columns: repeat(auto-fit, minmax(49%, 1fr));
			gap: 14px;
			.two {
				min-height: 52px;
			}
		}
	}
</style>
