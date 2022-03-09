<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  // tools
  import type { Meta } from "@/types";
  import { Env } from "@/tools/paths";
  import Utils from "@/tools/utils";
  // ui
  import Pagination from "@/ui/pagination.svelte";
  import Overlay from "@/ui/overlay.svelte";
  import Toolbar from "@/ui/toolbar.svelte";
  import SearchBar from "../../../ui/searchBar.svelte";
  // file
  import FileAdapter from "../adapter";
  import { By, Start, type File, type Params } from "../types";
  import CFile from "./file.svelte";
  import Uploader from "./uploader.svelte";

  /** add "select" option to selection overlay and dispatch event if this button clicked */
  export let withSelect: boolean = false;

  const dispatch = createEventDispatcher<{ selected: File }>();

  /** files loaded? */
  let loaded: boolean = false;
  /** file selected / overlay opened? */
  let toolsOverlay = false;
  /** selected file */
  let selected: {
    counter: number | null;
    file: File | null;
  } = { counter: null, file: null };
  /** response files */
  let files: Record<number, File> = {};
  /** response information */
  let meta: Meta;
  /** request params from portable mode */
  export let params: Params = undefined;
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

  /** get all files.
   * @param p request params
   *
   * @param withHistory set history params **by p parameter** after files loaded
   */
  async function getAll(p: Params = requestParams) {
    requestParams = requestParams;
    if (p.page < 1) {
      p.page = 1;
    }
    loaded = false;
    try {
      const result = await FileAdapter.getAll(p);
      files = result.data;
      meta = result.meta;
      loaded = true;
    } catch (err) {}
  }

  /** select file */
  function select(counter: number) {
    toolsOverlay = true;
    selected.counter = counter;
    selected.file = files[counter];
  }

  /** delete file */
  async function deleteFile(counter: number) {
    const isDelete = await window.$choose.confirm("delete file");
    if (!isDelete) {
      return;
    }
    toolsOverlay = false;
    try {
      const converted = getIDByCounter(counter);
      await FileAdapter.delete(converted);
      await deleteFromArray(counter);
    } catch (err) {}
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

  /** copy link to clipboard */
  async function copyLink(counter: number) {
    let message = "";
    const path = files[counter].path;
    const formattedPath = Env.getUploads() + "/" + path;
    try {
      await navigator.clipboard.writeText(formattedPath);
      message = "Link copied to clipboard.";
    } catch (err) {
      message = "Copy to clipboard: not have permission.";
    }
    window.$notify.add({ message });
    toolsOverlay = false;
  }

  /** play audio by url */
  function playAudio(url: URL) {
    window.$player.playlist = { position: 0, sources: [url.href] };
    window.$player.play();
  }

  /** delete file from files array */
  async function deleteFromArray(counter: number) {
    delete files[counter];
    files = files;
    await refresh();
  }

  /** set 'start' param and get files */
  function setStart(start: Start = Start.newest) {
    requestParams.start = start;
    requestParams.page = 1;
    getAll();
  }

  /** get file id by files counter id */
  function getIDByCounter(counter: number): string {
    return files[counter].id;
  }

  function search(val: string) {
    requestParams.filename = val;
    getAll();
  }
</script>

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
      {#each Object.entries(files) as [id, file]}
        <CFile {file} on:selected={() => select(parseInt(id, 10))} />
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

  {#if toolsOverlay}
    <Overlay onClose={() => (toolsOverlay = false)}>
      <div class="overlay">
        {#if withSelect}
          <div
            class="overlay__item"
            on:click={() => {
              dispatch("selected", selected.file);
            }}
          >
            select
          </div>
        {/if}
        {#if selected.file.extensionType.selected === "AUDIO"}
          <div
            class="overlay__item"
            on:click={() => playAudio(selected.file.pathConverted)}
          >
            play
          </div>
        {/if}
        <div class="overlay__item" on:click={() => copyLink(selected.counter)}>
          copy link
        </div>
        <div
          class="overlay__item"
          on:click={() => deleteFile(selected.counter)}
        >
          delete
        </div>
      </div>
    </Overlay>
  {/if}
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

  .overlay {
    width: 100%;
    &__item {
      height: 64px;
      width: 100%;
      font-size: 1rem;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    &__item:hover {
      background-color: var(--color-hover);
    }
  }
</style>
