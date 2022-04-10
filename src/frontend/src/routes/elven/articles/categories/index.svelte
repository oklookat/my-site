<script lang="ts">
	// ui
	import ToolbarBig from '$lib_elven/components/toolbar_big.svelte';
	// categories
	import type { Category } from '$lib_elven/types/articles/categories';
	import CompCategory from '$lib_elven/components/category.svelte';
	import CategoryNew from '$lib_elven/components/category_new.svelte';
	import type { Data } from '$lib_elven/types';
	import NetworkCategory from '$lib_elven/network/network_category';

	export let items: Data<Category>;

	/** create new cat? */
	let createNew = false;

	/** add or cancel adding cat? */
	let addText: 'add' | 'cancel' = 'add';
	$: createNew, onCreateNewChanged();
	const onCreateNewChanged = () => {
		addText = createNew ? 'cancel' : 'add';
	};

	function onAddClicked() {
		createNew = !createNew;
	}

	function onCancel() {
		createNew = false;
	}

	function getIDByCounter(counter: number): string | null {
		const _cat = items.data[counter];
		if (!_cat) {
			return null;
		}
		return _cat.id;
	}

	async function onCatAdded(name: string) {
		createNew = false;
		const cat: Category = {
			name: name
		};
		await addNew(cat);
	}

	async function addNew(cat: Category) {
		try {
			await NetworkCategory.create(cat);
		} catch (err) {}
	}

	async function deleteCat(counter: number) {
		const sure = await window.$choose.confirm('Delete category');
		if (!sure) {
			return;
		}
		const id = getIDByCounter(counter);
		if (!id) {
			return;
		}
		try {
			await NetworkCategory.delete(id);
			delete items.data[counter];
			items = items;
		} catch (err) {}
	}

	async function rename(counter: number, newName: string) {
		// store old and set new name
		const oldName = items.data[counter].name;
		items.data[counter].name = newName;
		//
		const cat: Category = {
			id: getIDByCounter(counter),
			name: newName
		};
		try {
			await NetworkCategory.rename(cat);
			items = items;
		} catch (err) {
			// revert name
			items.data[counter].name = oldName;
		}
	}
</script>

<svelte:head>
	<title>elven: article categories</title>
</svelte:head>

<div class="categories base__container">
	<ToolbarBig>
		<!-- prevent cancel block by textarea unfocus -->
		{#if !createNew}
			<div class="item" on:click={() => onAddClicked()}>
				{addText}
			</div>
		{/if}
	</ToolbarBig>

	{#if createNew}
		<CategoryNew on:submit={(e) => onCatAdded(e.detail)} on:cancel={() => onCancel()} />
	{/if}

	<div class="categories__list">
		{#each Object.entries(items.data) as [counter, cat]}
			<CompCategory
				{cat}
				on:delete={() => deleteCat(parseInt(counter, 10))}
				on:rename={(e) => rename(parseInt(counter, 10), e.detail)}
			/>
		{/each}
	</div>
</div>

<style lang="scss">
	.categories {
		&__list {
			display: flex;
			flex-direction: column;
			gap: 12px;
		}
	}
</style>
