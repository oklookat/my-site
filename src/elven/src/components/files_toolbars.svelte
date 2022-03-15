<script lang="ts">
    // tools
    import { Env } from "@/tools/paths";
    // ui
    import Overlay from "@/components/overlay.svelte";
    import Toolbar from "@/components/toolbar.svelte";
    import SearchBar from "@/components/search_bar.svelte";
    // files
    import type { Params } from "@/types/files";
    import { Start } from "@/types/files";
    import FilesUploader from "@/components/files_uploader.svelte";
    import { createEventDispatcher } from "svelte";

    export let params: Params;

    const dispatch = createEventDispatcher<{
        /** on request param changed */
        paramChanged: void;
    }>();

    /** set 'start' param */
    function setStart(start: Start = Start.newest) {
        params.start = start;
        params.page = 1;
        dispatch("paramChanged");
    }

    /** search by filename */
    function search(val: string) {
        params.filename = val;
        dispatch("paramChanged");
    }

    function onUploaded() {
        dispatch("paramChanged");
    }
</script>

<div class="toolbars">
    <FilesUploader on:uploaded={onUploaded} />

    <div class="oneline">
        <div class="sort">
            <Toolbar>
                <div class="sort-by-old">
                    {#if params.start === Start.newest}
                        <div
                            class="item"
                            on:click={() => setStart(Start.oldest)}
                        >
                            newest
                        </div>
                    {/if}
                    {#if params.start === Start.oldest}
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
