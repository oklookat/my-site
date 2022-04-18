<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	// ui
	import Select from '$lib_elven/components/select.svelte';
	// category
	import type { Category } from '$lib_elven/types/articles/categories';
	import NetworkCategory from '$lib_elven/network/network_category';
	import type { Counter, Items } from '$lib_elven/types';

	const dispatch = createEventDispatcher<{
		/** when new category selected */
		selected: { counter: Counter; category: Category };
	}>();

	/** currently selected category counter */
	export let selected: Counter | null = null;

	/** custom categories (adds before server categories) */
	export let customCategories: Record<Counter, Category> = undefined;

	/** categories from server */
	let serverCategories: Record<Counter, Category> = undefined;

	/** custom + server categories */
	let combinedCategories: Record<Counter, Category> = undefined;

	onMount(async () => {
		await getCategories();
		combine();
	});

	/** get categories from server */
	async function getCategories() {
		try {
			const networkCategory = new NetworkCategory('');
			const resp = await networkCategory.getAll();
			if (resp.status === 200) {
				const items: Items<Category> = await resp.json();
				serverCategories = items.data;
				return;
			}
		} catch (err) {}
		serverCategories = undefined;
	}

	/** categories for <select/> component */
	let selectable: Record<string, string> = {};

	/** combine categories */
	function combine() {
		combinedCategories = {};

		// combine cats
		if (customCategories) {
			combinedCategories = Object.assign({}, customCategories);
		}
		if (serverCategories) {
			combinedCategories = Object.assign(combinedCategories, serverCategories);
		}

		// make selectable from them
		selectable = {};
		for (const counter in combinedCategories) {
			const cat = combinedCategories[counter];
			selectable[counter] = cat.name;
		}
		selectable = selectable;
	}

	function onCategoryChanged(counter: string) {
		dispatch('selected', {
			counter: counter,
			category: combinedCategories[counter]
		});
	}
</script>

<Select bind:selectable bind:selected on:selected={(e) => onCategoryChanged(e.detail)} />
