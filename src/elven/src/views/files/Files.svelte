<script lang="ts">
  import { onMount } from "svelte";
  import FileAdapter from "@/common/adapters/Main/FileAdapter";
  import { IMeta, iMetaDefault } from "@/types/global";
  import type { IFile } from "@/types/file";
  import FilesList from "@/components/parts/FilesList.svelte";
  import Pagination from "@/components/ui/Pagination.svelte";

  let filesLoaded: boolean = false;
  // service
  let isSortOverlayActive: boolean = false;
  let sortBy: string = "created";
  let sortFirst: string = "newest";
  // meta
  let currentPage: number = 1;
  let totalPages: number = 1;
  // files
  let files: Array<IFile> = [];
  let meta: IMeta = iMetaDefault;
  let show: "newest";
  // file input for upload
  let input;

  onMount(() => {
    getFiles();
  });

  function getFiles(pageA = currentPage, showA = show) {
    currentPage = pageA;
    show = showA;
    filesLoaded = false;
    FileAdapter.getFiles(pageA, showA).then((result) => {
      files = result.data;
      meta = result.meta;
      currentPage = meta.current_page;
      totalPages = meta.total_pages;
      filesLoaded = true;
    });
  }

  async function refresh() {
    let isTrueFiles = filesLoaded && files.length < 1;
    // if (isTrueFiles) { // no files in current page
    //     while (isTrueFiles) {
    //         // moving back until the pages ends or data appears
    //         this.currentPage--
    //         await this.getFiles()
    //         if (this.currentPage <= 1) {
    //             break
    //         }
    //         isTrueFiles = this.loaded.value && this.files.value.length < 1
    //     }
    // }
  }

  async function deleteFile(file: IFile) {
    const isDelete = confirm("delete file?");
    if (isDelete) {
      FileAdapter.delete(file.id).then(() => {
        deleteFileFromArray(file);
      });
    }
  }

  async function onFileInputChange(event) {
    const files = event.target.files;
    if (files.length < 1) {
      return 0;
    }
    await FileAdapter.upload(files);
    await getFiles();
  }

  function onUploadClick() {
    if (!input.value) {
      return;
    }
    input.value.value = "";
    input.value.click();
  }

  function deleteFileFromArray(file: IFile) {
    const index = files.indexOf(file);
    files.splice(index, 1);
    return true;
  }
</script>

<main>
  <div class="files__container">
    <div class="files__tools">
      <div class="file__upload">
        <div class="file__upload-butt" on:click={() => onUploadClick()}>
          upload
        </div>
        <input
          bind:this={input}
          id="file__input"
          type="file"
          multiple
          style="display: none"
          on:input={onFileInputChange}
        />
      </div>
    </div>

    {#if filesLoaded && files.length < 1}
      <div class="files__404">
        <div>no files :(</div>
      </div>
    {/if}

    <FilesList :files="files" @delete="deleteFile($event)" />

    <Pagination
      :total-pages="totalPages"
      :current-page="currentPage"
      @changed="getFiles($event)"
    />
  </div>
</main>

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
