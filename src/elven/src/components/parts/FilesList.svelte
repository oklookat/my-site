<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { IFile } from "@/types/FileTypes";
  import { Env } from "@/tools/Paths";
  import Dates from "@/tools/Dates";
  import Sizes from "@/tools/Sizes";
  import Extensions from "@/tools/Extensions";
  import Overlay from "@/components/ui/Overlay.svelte";

  let selected: IFile | null = null;
  let selectionOverlay = false;

  export let files: Array<IFile> = [];
  let normalized: Array<IFile> = [];

  $: watchNormalized(files);
  function watchNormalized(files: Array<IFile>) {
    normalized = files.map((file) => {
      file = converter(file);
      return file;
    });
  }

  const dispatch = createEventDispatcher();

  function onSelected(file: IFile) {
    selectionOverlay = true;
    selected = file;
    dispatch("selected", file);
  }

  function onDelete(file: IFile) {
    selectionOverlay = false;
    dispatch("delete", file);
  }

  async function copyLink(file: IFile) {
    navigator.clipboard
      .writeText(file.path)
      .then(() => {
        selectionOverlay = false;
        window.$elvenNotify.add("Link copied to clipboard.");
      })
      .catch(() => {
        window.$elvenNotify.add("Copy to clipboard: not have permission.");
      });
  }

  // convert file path, extension etc
  function converter(file: IFile): IFile {
    file.path = `${Env.getUploads()}/${file.path}`;
    file.extensionType = Extensions.getType(file.extension);
    file.sizeConverted = Sizes.convert(file.size);
    file.createdAtConverted = Dates.convert(file.created_at);
    return file;
  }

  function playAudio(path: string) {
    window.$elvenPlayer.play(path);
  }
</script>

<div class="files__list">
  {#each normalized as file (file.id)}
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
  {/each}

  <Overlay
    bind:active={selectionOverlay}
    on:deactivated={() => {
      selectionOverlay = false;
      selected = null;
    }}
  >
    <div class="overlay__selected">
      {#if selected.extensionType === "audio"}
        <div
          class="overlay__item file__play"
          on:click={() => playAudio(selected.path)}
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

<style>
  .files__list {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
    min-height: 42px;
    gap: 12px;
  }

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
  }

  .file__item {
    font-size: 0.9rem;
    line-height: 1.5rem;
    margin-top: 8px;
    margin-left: 12px;
    margin-right: 12px;
  }

  .file__meta {
    display: flex;
    flex-direction: row;
    color: var(--color-text-inactive);
  }

  .file__main {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .file__name {
    font-size: 1.1rem;
    line-height: 2rem;
    letter-spacing: 0.0099rem;
  }

  .file__size {
    font-size: 1rem;
  }

  .overlay__selected {
    width: 100%;
  }

  .overlay__item {
    height: 64px;
    width: 100%;
    font-size: 1rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .overlay__item:hover {
    background-color: var(--color-hover);
  }
</style>
