<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  // tools
  import type { Meta } from "@/types";
  import { Env } from "@/tools/paths";
  import Utils from "@/tools";
  // ui
  import Pagination from "@/components/pagination.svelte";
  import Overlay from "@/components/overlay.svelte";
  import Toolbar from "@/components/toolbar.svelte";
  import SearchBar from "@/components/search_bar.svelte";
  // file
  import NetworkFile from "@/network/network_files";
  import { By, Start, type File, type Params } from "@/types/files";
  import CFile from "@/components/file.svelte";
  import Uploader from "@/components/files_uploader.svelte";
  import FileActions from "@/components/file_actions.svelte";

  /** add "select" option to selection overlay and dispatch event if this button clicked */
  export let withSelect: boolean = false;

  /** request params from portable mode */
  export let params: Params = undefined;

  const dispatch = createEventDispatcher<{ selected: File }>();

  /** files loaded? */
  let loaded = false;

  /** is file selected? */
  let isSelected = false;

  /** selected file */
  let selected: {
    counter: number | null;
    file: File | null;
    mouseEvent: MouseEvent;
  } = { counter: null, file: null, mouseEvent: null };

  /** response files */
  let files: Record<number, File> = {};

  /** response information */
  let meta: Meta;

  /** request params */
  let requestParams: Params = {
    page: 1,
    start: Start.newest,
    by: By.created,
  };

  onMount(() => {
    if (params) {
      Object.assign(requestParams, params);
    }
    getAll(undefined);
  });

  /** get all files */
  async function getAll(p: Params = requestParams) {
    requestParams = requestParams;
    if (p.page < 1) {
      p.page = 1;
    }
    loaded = false;
    try {
      const result = await NetworkFile.getAll(p);
      files = result.data;
      meta = result.meta;
      loaded = true;
    } catch (err) {}
  }

  /** select file */
  function select(file: File, mouseEvent: MouseEvent, counter: number) {
    selected.counter = counter;
    selected.mouseEvent = mouseEvent;
    selected.file = files[counter];
    isSelected = true;
  }

  /** refresh files */
  async function refresh() {
    while (true) {
      const filesLength = Utils.getObjectLength(files);
      const noFiles = loaded && filesLength < 1;
      if (!noFiles) {
        break;
      }
      requestParams.page--;
      try {
        await getAll();
      } catch (err) {
        break;
      }
      if (requestParams.page < 2) {
        break;
      }
    }
  }

  /** when page changed */
  function onPageChanged(page: number) {
    requestParams.page = page;
    getAll();
  }

  /** set 'start' param and get files */
  function setStart(start: Start = Start.newest) {
    requestParams.start = start;
    requestParams.page = 1;
    getAll();
  }

  /** search by filename */
  function search(val: string) {
    requestParams.filename = val;
    getAll();
  }

  /** on selected file deleted */
  function onDeleted() {
    isSelected = false;
    deleteFromArray(selected.counter);
  }

  /** delete file from files array */
  async function deleteFromArray(counter: number) {
    delete files[counter];
    files = files;
    await refresh();
  }
</script>

{#if isSelected}
  <FileActions
    {withSelect}
    file={selected.file}
    mouseEvent={selected.mouseEvent}
    onDisabled={() => (isSelected = false)}
    onDeleted={() => onDeleted()}
  />
{/if}

<div class="files base__container">
  <div class="toolbars">
    <Uploader onUploaded={() => getAll()} />

    <div class="oneline">
      <div class="sort">
        <Toolbar>
          <div class="sort-by-old">
            {#if requestParams.start === Start.newest}
              <div class="item" on:click={() => setStart(Start.oldest)}>
                newest
              </div>
            {/if}
            {#if requestParams.start === Start.oldest}
              <div class="item" on:click={() => setStart(Start.newest)}>
                oldest
              </div>
            {/if}
          </div>
        </Toolbar>
      </div>
      <div class="search">
        <SearchBar on:search={(e) => search(e.detail)} placeholder="search" />
      </div>
    </div>
  </div>

  <div class="list">
    {#if loaded && Utils.getObjectLength(files) > 0}
      {#each Object.entries(files) as [counter, file]}
        <CFile
          {file}
          onSelected={(file, event) =>
            select(file, event, parseInt(counter, 10))}
        />
      {/each}
    {/if}
  </div>

  <div class="pages">
    {#if loaded && meta && meta.total_pages && meta.current_page}
      <Pagination
        total={meta.total_pages}
        current={meta.current_page}
        on:changed={(e) => onPageChanged(e.detail)}
      />
    {/if}
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
