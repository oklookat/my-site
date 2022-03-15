<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  // tools
  import type { Meta } from "@/types";
  import Utils from "@/tools";
  // ui
  import Pagination from "@/components/pagination.svelte";
  // file
  import NetworkFile from "@/network/network_files";
  import { By, Start, type File, type Params } from "@/types/files";
  import CFile from "@/components/file.svelte";
  import FileActions from "@/components/file_actions.svelte";
  import FilesToolbars from "@/components/files_toolbars.svelte";

  /** add "select" option to selection overlay and dispatch event if this button clicked */
  export let withSelect: boolean = false;

  const dispatch = createEventDispatcher<{
    /** on 'select' option clicked on file */
    selected: File;
  }>();

  /** request params from portable mode */
  export let params: Params = undefined;

  /** files loaded? */
  let loaded = false;

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

  /** is file selected? */
  let isSelected = false;

  /** selected file */
  let selected: {
    counter: number | null;
    file: File | null;
    mouseEvent: MouseEvent;
  } = { counter: null, file: null, mouseEvent: null };

  onMount(() => {
    // set request params if portable mod active
    if (params) {
      Object.assign(requestParams, params);
    }
    getAll(undefined);
  });

  /** on request param changed */
  async function onParamChanged() {
    await getAll();
  }

  /** when page changed */
  function onPageChanged(page: number) {
    requestParams.page = page;
    getAll();
  }

  /** on selected file deleted */
  function onDeleted() {
    isSelected = false;
    deleteFromArray(selected.counter);
  }

  /** select file */
  function select(file: File, mouseEvent: MouseEvent, counter: number) {
    selected.counter = counter;
    selected.mouseEvent = mouseEvent;
    selected.file = files[counter];
    isSelected = true;
  }

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

  /** refresh files */
  async function refresh() {
    const getData = async () => {
      await getAll();
      return files;
    };
    const setPage = (val: number) => (requestParams.page = val);
    await Utils.refresh(requestParams.page, setPage, getData);
  }

  /** delete file from files array */
  async function deleteFromArray(counter: number) {
    delete files[counter];
    files = files;
    await refresh();
  }

  /** on 'select' button clicked on selected file */
  function onSelectClicked() {
    if (!withSelect) {
      return;
    }
    dispatch("selected", selected.file);
  }
</script>

{#if isSelected}
  <FileActions
    {withSelect}
    file={selected.file}
    mouseEvent={selected.mouseEvent}
    onDisabled={() => (isSelected = false)}
    onDeleted={() => onDeleted()}
    {onSelectClicked}
  />
{/if}

<div class="files base__container">
  <FilesToolbars
    bind:params={requestParams}
    on:paramChanged={async () => onParamChanged()}
  />

  <div class="list">
    {#if loaded && Utils.getRecordLength(files) > 0}
      {#each Object.entries(files) as [counter, file]}
        <CFile
          {file}
          onSelected={(e) => select(file, e, parseInt(counter, 10))}
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
