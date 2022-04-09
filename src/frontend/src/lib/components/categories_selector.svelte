<script lang="ts">
    import { createEventDispatcher, onMount } from "svelte";
    // ui
    import Select from "$lib/components/select.svelte";
    // category
    import type { Category } from "$lib/types/articles/categories";
import NetworkCategory from "$lib/network/network_category";

    /** initial id */
    export let selectedID: string | null | undefined = undefined;

    /** custom categories. Puts before server categories */
    export let customCategories: Record<number, Category> = {};

    const dispatch = createEventDispatcher<{
        /** dispatch new category when changed */
        changed: Category | null;
    }>();

    /** categories from server */
    let categories: Record<number, Category> = {};

    /** fetched categories for select element */
    let selectData: {
        value: string;
        text: string;
    }[] = [];

    /** currently selected category value (counter) */
    let selectedValue: string | undefined;

    onMount(async () => {
        await getCategories();
    });

    $: onSelectedIDChanged(selectedID);
    function onSelectedIDChanged(newVal?: string) {
        makeCategoriesSelectable();
    }

    async function getCategories() {
        // get categories
        try {
            const networkCategory = new NetworkCategory('')
            const items = await networkCategory.getAll()
            categories = items.data
        } catch (err) {
            return Promise.reject();
        }
        makeCategoriesSelectable();
        return Promise.resolve();
    }

    /** format categories from server and put result in catsSelectData */
    function makeCategoriesSelectable() {
        selectData = [];

        // format custom categories for select element
        fetchCategories(customCategories);

        // create "no category" item
        const noCategory = { value: "", text: "No category" };
        selectData.push(noCategory)

        // format categories for select element
        fetchCategories(categories);
    }

    /** format & add categories to select */
    function fetchCategories(items: Record<number, Category>) {
        // if without selected category id - set "no category" as default
        const initialID = selectedID;
        if (!initialID) {
            selectedValue = "";
        }
        // format custom categories for select element
        for (const [counter, _category] of Object.entries(items)) {
            const option = {
                value: counter,
                text: _category.name,
            };
            // same categories?
            const sameCats = initialID === _category.id;
            if (sameCats) {
                selectedValue = counter;
            }
            selectData.push(option);
        }
        // render
        selectData = selectData;

    }

    /** when category on select element changed */
    function onCategoryChanged(counter?: string) {
        // no counter = no category
        let newCat: Category | null = null;
        if (counter) {
            newCat = getCategoryByCounter(counter);
        }
        dispatch("changed", newCat);
    }

    /** get category by categories counter */
    function getCategoryByCounter(counter?: string | number): Category | null {
        let cat: Category | null = null;
        if (!counter) {
            return cat;
        }
        try {
            const isString = typeof counter === "string";
            const counterInt = isString ? parseInt(counter, 10) : counter;
            cat =
                categories[counterInt] || customCategories[counterInt] || null;
        } catch (err) {}
        return cat;
    }
</script>

<Select
    bind:options={selectData}
    bind:selected={selectedValue}
    on:selected={(e) => onCategoryChanged(e.detail)}
/>