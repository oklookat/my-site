<script lang="ts">
  import { onMount } from "svelte";
  import FileAdapter from "@/common/adapters/Main/FileAdapter";
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
  let input;

  onMount(() => {
    getFiles();
  });

  function getFiles(pageA = currentPage, showA = show) {
    currentPage = pageA;
    show = showA;
    isLoaded = false;
    FileAdapter.getFiles(pageA, showA).then((result) => {
      files = result.data;
      meta = result.meta;
      currentPage = meta.current_page;
      totalPages = meta.total_pages;
      perPage = meta.per_page;
      isLoaded = true;
    });
  }

  async function deleteFile(file: IFile) {
    const isDelete = confirm("delete file?");
    if (isDelete) {
      FileAdapter.delete(file.id).then(() => {
        deleteFromArray(file);
        refresh();
      });
    }
  }

  function deleteFromArray(file: IFile) {
    files = files.filter((f) => f !== file);
    if (files.length < perPage) {
      getFiles();
    }
  }

  async function refresh() {
    let isTrueFiles = isLoaded && files.length < 1;
    if (isTrueFiles) {
      // no files in current page
      while (isTrueFiles) {
        // moving back until the pages ends or data appears
        currentPage--;
        await getFiles();
        if (currentPage <= 1) {
          break;
        }
        isTrueFiles = isLoaded && files.length < 1;
      }
    }
  }

  async function onInputChange(event) {
    const files = event.target.files;
    if (files.length < 1) {
      return 0;
    }
    await FileAdapter.upload(files);
    await getFiles();
  }

  function onUploadClick() {
    if (!input) {
      return;
    }
    input.value = "";
    input.click();
  }
</script>

<div class="files__container">
  <div class="files__tools">
    <div class="file__upload">
      <div class="file__upload-butt" on:click={() => onUploadClick()}>
        upload
      </div>
      <input
        id="file__input"
        type="file"
        multiple
        style="display: none"
        bind:this={input}
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

<style scoped>
  .files__container {
    width: 95%;
    max-width: 512px;
    margin: auto;
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .files__tools {
    background-color: var(--color-level-1);
    font-size: 1rem;
    width: 100%;
    height: 54px;
    border-radius: var(--border-radius);
    display: flex;
    flex-direction: row;
    align-items: center;
  }

  .file__upload {
    margin-left: 12px;
  }

  .file__upload-butt {
    cursor: pointer;
  }

  .files__404 {
    background-color: var(--color-level-1);
    height: 240px;
    border-radius: var(--border-radius);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 24px;
  }
</style>
