<script lang="ts">
    import { PathTools } from "$lib/tools/paths";
    import Dates from "$lib/tools/dates";
    import Size from "$lib/tools/size";
    import Extension from "$lib/tools/extension";
    import type { File } from "$lib/types/files";

    export let file: File;

    /** when clicked on file */
    export let onSelected: (e: MouseEvent) => void;

    $: convert(file)

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

<div class="file base__card" on:click={onSelected}>
    <div class="meta">
        <div class="meta__item">
            {file.createdAtConverted}
        </div>
        <div class="meta__item">{file.sizeConverted}</div>
    </div>
    <div class="main">
        {#if file.extensionsSelector && file.extensionsSelector.selected === "IMAGE"}
            <div class="file__preview">
                <img
                    decoding="async"
                    loading="lazy"
                    src={file.pathConverted.href}
                    alt=""
                />
            </div>
        {:else if file.extensionsSelector && file.extensionsSelector.selected === "VIDEO"}
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
