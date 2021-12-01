<script lang="ts">
  import { onMount } from "svelte";
  import FileAdapter from "@/adapters/FileAdapter";
  import type { Meta } from "@/types/global";
  import type { File } from "@/types/file";
  import CFile from "@/components/File.svelte";
  import Pagination from "@/ui/Pagination.svelte";
  import Overlay from "@/ui/Overlay.svelte";
  import ToolbarBig from "@/ui/ToolbarBig.svelte"

  let isLoaded: boolean = false;
  // service
  let isSortOverlayActive: boolean = false;
  let sortBy: string = "created";
  let sortFirst: string = "newest";
  let selected: File | null = null;
  let selectionOverlay = false;
  // files
  let files: Array<File> = [];
  let meta: Meta;
  let show: "newest";
  let perPage: number = 0;
  let currentPage: number = 1;
  let totalPages: number = 1;
  // file input for upload
  let inputEL: HTMLInputElement;

  onMount(() => {
    getFiles();
  });

  async function getFiles(pageA = currentPage, showA = show) {
    if (pageA < 1) {
      pageA = 1;
    }
    currentPage = pageA;
    show = showA;
    isLoaded = false;
    try {
      const result = await FileAdapter.getAll(pageA, showA);
      files = result.data;
      meta = result.meta;
      currentPage = meta.current_page;
      totalPages = meta.total_pages;
      perPage = meta.per_page;
      isLoaded = true;
    } catch (err) {}
  }

  async function deleteFile(file: File) {
    const isDelete = confirm("Delete file?");
    if (!isDelete) {
      return;
    }
    try {
      await FileAdapter.delete(file.id);
      await deleteFromArray(file);
    } catch (err) {}
  }

  async function deleteFromArray(file: File) {
    files = files.filter((f) => f !== file);
    await refresh();
  }

  async function refresh() {
    let noFiles = isLoaded && (files.length < 1 || files.length < perPage);
    // no files in current page
    while (noFiles && currentPage > 1) {
      // moving back until the pages ends or data appears
      currentPage--;
      await getFiles();
      noFiles = isLoaded && files.length < 1;
    }
  }

  async function onInputChange(e) {
    const files = <FileList>e.target.files;
    if (files.length < 1) {
      return 0;
    }
    FileAdapter.upload(files).then(() => {
      getFiles();
    });
  }

  function onUploadClick() {
    if (!inputEL) {
      return;
    }
    inputEL.value = "";
    inputEL.click();
  }

  function onSelected(file: File) {
    selectionOverlay = true;
    selected = file;
  }

  function onDelete(file: File) {
    selectionOverlay = false;
    deleteFile(file);
  }

  async function copyLink(file: File) {
    try {
      await navigator.clipboard.writeText(file.path);
      selectionOverlay = false;
      const message = "Link copied to clipboard.";
      window.$elvenNotify.add({ message });
    } catch (err) {
      const message = "Copy to clipboard: not have permission.";
      window.$elvenNotify.add({ message });
    }
  }

  function playAudio(url: URL) {
    window.$elvenPlayer.playlist = { position: 0, sources: [url.href] };
    window.$elvenPlayer.play();
  }
</script>

<div class="files">
  <div class="files__toolbar">
    <ToolbarBig>
      <div class="files__upload" on:click={() => onUploadClick()}>
        upload
        <input
        id="file__input"
        type="file"
        multiple
        style="display: none"
        bind:this={inputEL}
        on:input={onInputChange}
      />
      </div>
    </ToolbarBig>
  </div>

  {#if isLoaded && files.length < 1}
    <div class="files__404">
      <div>No files :(</div>
    </div>
  {/if}

  {#if isLoaded && files.length > 0}
    <div class="files__list">
      {#each files as file (file.id)}
        <CFile {file} on:selected={(e) => onSelected(e.detail)} />
      {/each}

      <Overlay
        bind:active={selectionOverlay}
        on:deactivated={() => {
          selectionOverlay = false;
          selected = null;
        }}
      >
        <div class="overlay">
          {#if selected.extensionType === "audio"}
            <div
              class="overlay__item file__play"
              on:click={() => playAudio(selected.pathConverted)}
            >
              play
            </div>
          {/if}
          <div
            class="overlay__item file__copy-link"
            on:click={() => copyLink(selected)}
          >
            copy link
          </div>
          <div
            class="overlay__item file__delete"
            on:click={() => onDelete(selected)}
          >
            delete
          </div>
        </div>
      </Overlay>
    </div>

    <Pagination
      total={totalPages}
      current={currentPage}
      on:changed={(e) => getFiles(e.detail)}
    />
  {/if}
</div>

<style lang="scss">
  .files {
    width: 95%;
    max-width: 512px;
    margin: auto;
    display: flex;
    flex-direction: column;
    gap: 14px;
    &__404 {
      background-color: var(--color-level-1);
      height: 240px;
      border-radius: var(--border-radius);
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 24px;
    }
    &__upload {
      cursor: pointer;
    }
    &__list {
      height: 100%;
      width: 100%;
      display: flex;
      flex-direction: column;
      min-height: 42px;
      gap: 12px;
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
