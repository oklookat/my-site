<script lang="ts">
    import { createEventDispatcher } from "svelte";
    // ui
    import Toolbar from "@/components/toolbar.svelte";
    import ToolbarBig from "@/components/toolbar_big.svelte";
    // article
    import { Show, By, Start } from "@/types/articles";
    import type { Params } from "@/types/articles";
    import type { Category } from "@/types/articles/categories";
    import CategoriesSelector from "@/components/categories_selector.svelte";

    export let params: Params;

    const dispatch = createEventDispatcher<{
        /** on request param changed */
        paramChanged: void;
    }>();

    /** custom categories for searching by category */
    const customCategories: Record<number, Category> = {
        "-1": {
            id: "-1",
            name: "All",
        },
    };

    /** set 'by' param and get articles */
    function setBy(by: By = By.published) {
        params.by = by;
        params.page = 1;
        dispatch("paramChanged");
    }

    /** set 'start' param and get articles */
    function setStart(start: Start = Start.newest) {
        params.start = start;
        params.page = 1;
        dispatch("paramChanged");
    }

    /** set 'show' param and get articles */
    function setShow(show: Show) {
        params.show = show;
        params.page = 1;
        dispatch("paramChanged");
    }

    /** sort by category */
    function onCategoryChanged(cat: Category | null) {
        params.page = 1;
        // no categories
        params.without_category = cat === null;
        let catName = null;
        // not all categories / not 'No category'
        const isRealCategory = !!(cat && cat["name"] && cat.id !== "-1");
        if (isRealCategory) {
            catName = cat.name;
        }
        params.category_name = catName;
        dispatch("paramChanged");
    }
</script>

<div class="toolbars">
    <ToolbarBig>
        <a href="#/articles/create">new</a>
        <a href="#/articles/cats">categories</a>
    </ToolbarBig>

    <Toolbar>
        <CategoriesSelector
            {customCategories}
            selectedID={"-1"}
            on:changed={(e) => onCategoryChanged(e.detail)}
        />
    </Toolbar>

    <Toolbar>
        <div>
            {#if params.show === Show.published}
                <div class="pointer" on:click={() => setShow(Show.drafts)}>
                    published
                </div>
            {/if}
            {#if params.show === Show.drafts}
                <div class="pointer" on:click={() => setShow(Show.published)}>
                    drafts
                </div>
            {/if}
        </div>
        <div>
            {#if params.start === Start.newest}
                <div class="pointer" on:click={() => setStart(Start.oldest)}>
                    newest
                </div>
            {/if}
            {#if params.start === Start.oldest}
                <div class="pointer" on:click={() => setStart(Start.newest)}>
                    oldest
                </div>
            {/if}
        </div>
        <div>
            {#if params.by === By.updated}
                <div class="pointer" on:click={() => setBy(By.published)}>
                    by updated date
                </div>
            {/if}
            {#if params.by === By.published}
                <div class="pointer" on:click={() => setBy(By.created)}>
                    by published date
                </div>
            {/if}
            {#if params.by === By.created}
                <div class="pointer" on:click={() => setBy(By.updated)}>
                    by created date
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
    }
</style>
