<script lang="ts">
  import { onMount } from "svelte";
  import FileAdapter from "@/adapters/FileAdapter";
  import { IMeta, iMetaDefault } from "@/types/global";
  import type { IFile } from "@/types/file";
  import FilesList from "@/components/parts/FilesList.svelte";
  import Pagination from "@/components/ui/Pagination.svelte";

  let isLoaded: boolean = false;
  // service
  let isSortOverlayActive: boolean = false;
  let sortBy: string = "created";
  let sortFirst: string = "newest";
  // files
  let files: Array<IFile> = [];
  let meta: IMeta = iMetaDefault;
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
    if(pageA < 1) {
      pageA = 1
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

  async function deleteFile(file: IFile) {
    const isDelete = confirm("Delete file?");
    if (!isDelete) {
      return;
    }
    try {
      await FileAdapter.delete(file.id);
      await deleteFromArray(file);
    } catch (err) {}
  }

  async function deleteFromArray(file: IFile) {
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
</script>

<div class="files">
  <div class="files__tools">
    <div class="files__upload">
      <div class="files__upload-button" on:click={() => onUploadClick()}>
        upload
      </div>
      <input
        id="file__input"
        type="file"
        multiple
        style="display: none"
        bind:this={inputEL}
        on:input={onInputChange}
      />
    </div>
  </div>

  {#if isLoaded && files.length < 1}
    <div class="files__404">
      <div>No files :(</div>
    </div>
  {/if}

  {#if isLoaded && files.length > 0}
    <FilesList {files} on:delete={(e) => deleteFile(e.detail)} />
    <Pagination
      {totalPages}
      {currentPage}
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
    &__tools {
      background-color: var(--color-level-1);
      font-size: 1rem;
      width: 100%;
      height: 54px;
      border-radius: var(--border-radius);
      display: flex;
      flex-direction: row;
      align-items: center;
    }
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
      margin-left: 12px;
      &-button {
        cursor: pointer;
      }
    }
  }
</style>
