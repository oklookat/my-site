<script lang="ts">
    import { createEventDispatcher, onDestroy, onMount } from "svelte";
    // ui
    import Animation from "$lib/tools/animation";
    // files
    import type { File, Params } from "$lib/types/files";
    import Files from "$lib/views/files/index.svelte";

    const dispatch = createEventDispatcher<{ 
        /** on 'select' option clicked on file */
        selected: File, 
        /** on files closed */
        closed: void  }>();

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

    function onClosed() {
        dispatch("closed")
    }
</script>

<div class="overlay base__overlay" bind:this={container}>
    <div class="overlay__main">
        <div class="close pointer" on:click={onClosed}>
            <svg viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg"
                ><path
                    d="M289.94,256l95-95A24,24,0,0,0,351,127l-95,95-95-95A24,24,0,0,0,127,161l95,95-95,95A24,24,0,1,0,161,385l95-95,95,95A24,24,0,0,0,385,351Z"
                /></svg
            >
        </div>
        <div class="files">
            <Files
                isSelectMode={true}
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
                svg {
                    fill: red;
                    width: 30px;
                    height: 30px;
                }
            }
            .files {
                margin-bottom: 32px;
            }
        }
    }
</style>
