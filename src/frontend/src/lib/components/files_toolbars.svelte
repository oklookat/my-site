<script lang="ts">
    // tools
    // ui
    import Toolbar from "$lib/components/toolbar.svelte";
    import SearchBar from "$lib/components/search_bar.svelte";
    // files
    import type { Params } from "$lib/types/files";
    import { Start } from "$lib/types/files";
    import FilesUploader from "$lib/components/files_uploader.svelte";
    import { createEventDispatcher } from "svelte";

    export let params: Params;

    const dispatch = createEventDispatcher<{
        /** on request param changed */
        paramChanged: {name: string, val: string};
        upload: void
    }>();

    /** set 'start' param */
    function setStart(start: Start = Start.newest) {
        params.start = start;
        dispatch("paramChanged", {name: 'start', val: start});
    }

    /** search by filename */
    function search(val: string) {
        params.filename = val;
        dispatch("paramChanged", {name: 'filename', val: val});
    }

    function onUpload() {
        dispatch("upload");
    }
</script>

<div class="toolbars">
    <FilesUploader on:uploaded={onUpload} />

    <div class="oneline">
        <div class="sort">
            <Toolbar>
                <div class="sort-by-old">
                    {#if params && params.start === Start.newest}
                        <div
                            class="item"
                            on:click={() => setStart(Start.oldest)}
                        >
                            newest
                        </div>
                    {/if}
                    {#if params && params.start === Start.oldest}
                        <div
                            class="item"
                            on:click={() => setStart(Start.newest)}
                        >
                            oldest
                        </div>
                    {/if}
                </div>
            </Toolbar>
        </div>
        <div class="search">
            <SearchBar
                on:search={(e) => search(e.detail)}
                placeholder="search"
            />
        </div>
    </div>
</div>

<style lang="scss">
    .toolbars {
        display: flex;
        flex-direction: column;
        gap: 12px;
        width: 100%;
        .oneline {
            display: flex;
            gap: 14px;
            width: 100%;
            .sort {
                width: 50%;
            }
            .search {
                height: 54px;
                width: 50%;
            }
        }
    }
</style>
