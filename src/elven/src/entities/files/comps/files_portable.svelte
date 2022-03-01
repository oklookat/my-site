<script lang="ts">
    import { createEventDispatcher, onDestroy, onMount } from "svelte";
    // ui
    import Animation from "../../../tools/animation";
    // files
    import type { File, Params } from "../types";
    import Files from "./index.svelte";

    const dispatch = createEventDispatcher<{ selected: File }>();

    export let onClose: () => void;
    export let params: Params = undefined;

    let container: HTMLDivElement;

    onMount(() => {
        document.body.classList.add("no-scroll");
        Animation.fadeIn(container, 10);
    });

    onDestroy(() => {
        document.body.classList.remove("no-scroll");
    });

    function onSelected(file: File) {
        dispatch("selected", file);
    }
</script>

<div class="overlay overlay-foundation" bind:this={container}>
    <div class="overlay__main">
        <div class="close pointer" on:click={onClose}>close</div>
        <div class="files">
            <Files
                withSelect={true}
                {params}
                on:selected={(e) => onSelected(e.detail)}
            />
        </div>
    </div>
</div>

<style lang="scss">
    .overlay {
        background-color: var(--color-body);
        display: block;
        &__main {
            z-index: 9999;
            overflow: auto;
            width: 100%;
            height: 100%;
            display: grid;
            grid-template-rows: max-content;
            grid-template-columns: 1fr;
            gap: 14px;
            .close {
                width: 100%;
                height: 48px;
                background-color: var(--color-level-1);
                display: flex;
                justify-content: center;
                align-items: center;
            }
            .files {
                margin-bottom: 32px;
            }
        }
    }
</style>
