<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { TFile } from "@/types/file";
    import { Env } from "@/tools/Paths";
    import Dates from "@/tools/Dates";
    import Sizes from "@/tools/Sizes";
    import Extensions from "@/tools/Extensions";

    const dispatch = createEventDispatcher<{ selected: TFile }>();

    export let file: TFile;
    $: watchFile(file);

    function watchFile(file: TFile) {
        file = converter(file);
    }

    function onSelected(file: TFile) {
        dispatch("selected", file);
    }

    // convert file path, extension etc
    function converter(file: TFile): TFile {
        file.path = `${Env.getUploads()}/${file.path}`;
        file.extensionType = Extensions.getType(file.extension);
        file.sizeConverted = Sizes.convert(file.size);
        file.createdAtConverted = Dates.convert(file.created_at);
        return file;
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
                <img src={file.path} alt="" />
            </div>
        {:else if file.extensionType === "video"}
            <div class="file__item file__preview" on:click|stopPropagation>
                <video controls src={file.path}>
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
        height: 100%;
        display: flex;
        flex-direction: column;
        padding-bottom: 12px;
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
            display: flex;
            flex-direction: column;
            gap: 4px;
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
