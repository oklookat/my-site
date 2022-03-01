<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { File } from "../types";
    import { PathTools } from "@/tools/paths";
    import Dates from "@/tools/dates";
    import Size from "@/tools/size";
    import Extension from "@/tools/extension";

    const dispatch = createEventDispatcher<{ selected: File }>();

    export let file: File;
    $: watchFile(file);

    function watchFile(file: File) {
        convert(file);
    }

    function onSelected(file: File) {
        dispatch("selected", file);
    }

    /** convert file path, extension etc */
    // TODO: split converter functions(?)
    function convert(file: File) {
        let isNeedPath = !(file.pathConverted instanceof URL);
        if (isNeedPath) {
            file.pathConverted = PathTools.getUploadsWith(file.path);
        }
        if (!file.extensionType) {
            file.extensionType = Extension.getType(file.extension);
        }
        if (!file.sizeConverted) {
            file.sizeConverted = Size.convert(file.size);
        }
        if (!file.createdAtConverted) {
            file.createdAtConverted = Dates.convert(file.created_at);
        }
    }
</script>

<div class="file" on:click={() => onSelected(file)}>
    <div class="file__meta">
        <div class="file__item file__uploaded-date">
            {file.createdAtConverted}
        </div>
    </div>
    <div class="file__main">
        {#if file.extensionType === "image"}
            <div class="file__item file__preview" on:click|stopPropagation>
                <img decoding="async" loading="lazy" src={file.pathConverted.href} alt="" />
            </div>
        {:else if file.extensionType === "video"}
            <div class="file__item file__preview" on:click|stopPropagation>
                <video controls src={file.pathConverted.href}>
                    <track default kind="captions" srclang="en" src="" />
                </video>
            </div>
        {/if}
        <div class="file__item file__name">{file.original_name}</div>
        <div class="file__item file__size">{file.sizeConverted}</div>
    </div>
</div>

<style lang="scss">
    .file {
        box-shadow: 0 0 41px 0 rgba(0, 0, 0, 0.05);
        min-height: 42px;
        border-radius: var(--border-radius);
        background-color: var(--color-level-1);
        cursor: pointer;
        width: 100%;
        height: fit-content;
        padding-bottom: 12px;
        display: grid;
        grid-template-columns: 1fr;
        grid-template-rows: 1fr;
        gap: 8px;
        &__item {
            font-size: 0.9rem;
            line-height: 1.5rem;
            margin-top: 8px;
            margin-left: 12px;
            margin-right: 12px;
        }
        &__meta {
            display: flex;
            flex-direction: row;
            color: var(--color-text-inactive);
        }
        &__main {
            height:auto;
            display: flex;
            flex-direction: column;
            gap: 4px;
        }
        &__preview {
            max-width: 100%;
        }
        &__name {
            font-size: 1.1rem;
            line-height: 2rem;
            letter-spacing: 0.0099rem;
        }
        &__size {
            font-size: 1rem;
        }
    }
</style>
