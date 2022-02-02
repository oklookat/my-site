<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { Category } from "../types";
    import New from "./new.svelte";

    export let cat: Category;

    const dispatch = createEventDispatcher<{
        delete: void;
        rename: string;
    }>();

    let renameMode = false;

    function onDeleteClicked() {
        dispatch("delete");
    }

    function onRenameClicked() {
        renameMode = true;
    }

    function onFocusOut() {
        renameMode = false;
    }

    function onSubmit(newData: string) {
        renameMode = false;
        if(cat.name === newData) {
            return
        }
        dispatch("rename", newData);
    }
</script>

{#if renameMode}
    <div class="category">
        <New
            data={cat.name}
            on:cancel={() => onFocusOut()}
            on:submit={(e) => onSubmit(e.detail)}
        />
    </div>
{/if}

{#if !renameMode}
    <div class="category pointer" on:click|stopPropagation={onRenameClicked}>
        <div class="category__name">
            {cat.name}
        </div>
        <div
            class="category__delete pointer"
            on:click|stopPropagation={onDeleteClicked}
        >
            <svg
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
                xmlns:xlink="http://www.w3.org/1999/xlink"
                x="0px"
                y="0px"
                viewBox="0 0 1792 1792"
                style="enable-background:new 0 0 1792 1792;"
                xml:space="preserve"
            >
                <path
                    d="M1082.2,896.6l410.2-410c51.5-51.5,51.5-134.6,0-186.1s-134.6-51.5-186.1,0l-410.2,410L486,300.4
	c-51.5-51.5-134.6-51.5-186.1,0s-51.5,134.6,0,186.1l410.2,410l-410.2,410c-51.5,51.5-51.5,134.6,0,186.1
	c51.6,51.5,135,51.5,186.1,0l410.2-410l410.2,410c51.5,51.5,134.6,51.5,186.1,0c51.1-51.5,51.1-134.6-0.5-186.2L1082.2,896.6z"
                />
            </svg>
        </div>
    </div>
{/if}

<style lang="scss">
    .category {
        background-color: var(--color-level-1);
        font-size: 1.3rem;
        border-radius: 4px;
        display: flex;
        align-items: center;
        width: 100%;
        height: 100%;
        min-height: 48px;
        align-content: center;
        &__name {
            padding: 8px;
        }
        &__delete {
            width: max-content;
            margin-left: auto;
            margin-right: 6px;
            height: 100%;
            display: flex;
            justify-content: center;
            svg {
                fill: red;
                width: 20px;
                height: 20px;
            }
        }
    }
</style>
