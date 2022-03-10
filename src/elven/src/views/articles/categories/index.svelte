<script lang="ts">
    import { onMount } from "svelte";
    // import type { Meta } from "@/types";
    // ui
    import ToolbarBig from "@/components/toolbar_big.svelte";
    // categories
    import type { Category } from "@/types/articles/categories";
    import CompCategory from "@/components/category.svelte";
    import NetworkCateogries from "@/network/network_categories";
    import CategoryNew from "@/components/category_new.svelte";

    /** response cats */
    let cats: Record<number, Category> = {};
    /** response information */
    // let meta: Meta;
    /** content loaded? */
    let loaded = false;
    /** add or cancel adding cat? */
    let addText: "add" | "cancel" = "add";
    /** create new cat? */
    let createNew = false;
    $: createNew, onCreateNewChanged();
    const onCreateNewChanged = () => {
        addText = createNew ? "cancel" : "add";
    };

    onMount(() => {
        getAll();
    });

    function onAddClicked() {
        createNew = !createNew;
    }

    function onCancel() {
        createNew = false;
    }

    function getIDByCounter(counter: number): string | null {
        const _cat = cats[counter];
        if (!_cat) {
            return null;
        }
        return _cat.id;
    }

    async function onCatAdded(name: string) {
        createNew = false;
        const cat: Category = {
            name: name,
        };
        await addNew(cat);
    }

    async function getAll() {
        try {
            const result = await NetworkCateogries.getAll();
            cats = result.data;
            // meta = result.meta;
            loaded = true;
        } catch (err) {
            loaded = false;
        }
    }

    async function addNew(cat: Category) {
        try {
            await NetworkCateogries.create(cat);
            await getAll();
        } catch (err) {}
    }

    async function deleteCat(counter: number) {
        const sure = await window.$choose.confirm("delete category")
        if (!sure) {
            return;
        }
        const id = getIDByCounter(counter);
        if (!id) {
            return;
        }
        try {
            await NetworkCateogries.delete(id);
            delete cats[counter];
            cats = cats;
        } catch (err) {}
    }

    async function rename(counter: number, newName: string) {
        // store old and set new name
        const oldName = cats[counter].name;
        cats[counter].name = newName;
        //
        const cat: Category = {
            id: getIDByCounter(counter),
            name: newName,
        };
        try {
            await NetworkCateogries.rename(cat);
            cats = cats;
        } catch (err) {
            // revert name
            cats[counter].name = oldName;
        }
    }
</script>

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
        <CategoryNew
            on:submit={(e) => onCatAdded(e.detail)}
            on:cancel={() => onCancel()}
        />
    {/if}

    {#if loaded}
        <div class="categories__list">
            {#each Object.entries(cats) as [counter, cat]}
                <CompCategory
                    {cat}
                    on:delete={() => deleteCat(parseInt(counter, 10))}
                    on:rename={(e) => rename(parseInt(counter, 10), e.detail)}
                />
            {/each}
        </div>
    {/if}
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
