<script lang="ts">
    import { createEventDispatcher } from "svelte";

    import { PathTools } from "@/tools/paths";
    import Dates from "@/tools/dates";
    import Size from "@/tools/size";
    import Extension from "@/tools/extension";
    import type { File } from "@/types/files";
    // ui
    import Popup from "@/components/popup.svelte";
    let popupMouseEvent: MouseEvent;

    const dispatch = createEventDispatcher<{ selected: File }>();

    export let file: File;
    $: watchFile(file);

    function watchFile(file: File) {
        convert(file);
    }

    let isSelected = false;
    function onSelected(file: File, e: MouseEvent) {
        isSelected = true;
        popupMouseEvent = e;
        //dispatch("selected", file);
    }

    /** convert file path, extension etc */
    // TODO: split converter functions(?)
    function convert(file: File) {
        let isNeedPath = !(file.pathConverted instanceof URL);
        if (isNeedPath) {
            file.pathConverted = PathTools.getUploadsWith(file.path);
        }
        if (!file.extensionsSelector) {
            file.extensionsSelector = Extension.getSelector(file.extension);
        }
        if (!file.sizeConverted) {
            file.sizeConverted = Size.convert(file.size);
        }
        if (!file.createdAtConverted) {
            file.createdAtConverted = Dates.convert(file.created_at);
        }
    }
</script>

<!-- TODO: OVERLAY ON MOBILE / POPUP ON PC
TODO: FILE MANIPULATIONS DO IN HERE, NOT IN FILES LIST
TODO: AND WITH ARTICLE SAME -->
{#if isSelected}
    <Popup
        bind:mouseEvent={popupMouseEvent}
        onDisabled={() => (isSelected = false)}
    />
{/if}

<div class="file base__card" on:click={(e) => onSelected(file, e)}>
    <div class="meta">
        <div class="meta__item">
            {file.createdAtConverted}
        </div>
        <div class="meta__item">{file.sizeConverted}</div>
    </div>
    <div class="main">
        {#if file.extensionsSelector.selected === "IMAGE"}
            <div class="file__preview">
                <img
                    decoding="async"
                    loading="lazy"
                    src={file.pathConverted.href}
                    alt=""
                />
            </div>
        {:else if file.extensionsSelector.selected === "VIDEO"}
            <div class="file__preview" on:click|stopPropagation>
                <video controls src={file.pathConverted.href}>
                    <track default kind="captions" srclang="en" src="" />
                </video>
            </div>
        {/if}
        <div class="title">{file.original_name}</div>
    </div>
</div>

<style lang="scss">
    .file {
        &__preview {
            :global(img),
            :global(video) {
                width: 100%;
                max-height: 224px;
                object-fit: cover;
            }
        }
    }
</style>
