<template>
  <div class="files__container">
    <div class="files__tools">
      <div class="file__upload">
        <div style="cursor: pointer;" class="file__upload-butt" v-on:click="onUploadClick()">upload</div>
        <input
          ref="input"
          id="file__input"
          type="file"
          multiple
          style="display: none"
          @input="onFileInputChange($event)"
        />
      </div>
    </div>

    <div class="files__list" v-if="filesLoaded && files.length > 0">
      <div class="file" v-for="file in files" :key="file.id" v-on:click="select(file)">
        <div class="file__meta">
          <div class="file__item file__uploaded-date">{{ convertDateWrap(file.created_at) }}</div>
        </div>
        <div class="file__main">
          <div
            class="file__item file__preview"
            v-on:click.stop
            v-if="readableExtensionWrap(file.extension) === 'IMAGE'"
          >
            <img
              :src="convertPreviewPath(file.path)"
              v-if="readableExtensionWrap(file.extension) === 'IMAGE'"
            />
          </div>
          <div
            class="file__item file__preview"
            v-on:click.stop
            v-if="readableExtensionWrap(file.extension) === 'VIDEO'"
          >
            <video controls :src="convertPreviewPath(file.path)"></video>
          </div>

          <div class="file__item file__name">{{ file.original_name }}</div>
          <div class="file__item file__size">{{ convertSizeWrap(file.size) }}</div>
        </div>
      </div>
    </div>

    <div class="files__404" v-if="filesLoaded && files.length < 1">
      <div>no files :(</div>
    </div>

    <Overlay v-bind:active="isToolsOverlayActive" v-on:deactivated="isToolsOverlayActive = false">
      <div class="overlay__file-tools">
        <div
          class="overlay__item file__play"
          v-if="readableExtensionWrap(selectedFile.extension) === 'AUDIO'"
          v-on:click="playAudio(selectedFile)"
        >play</div>
        <div class="overlay__item file__copy-link" v-on:click="copyLink(selectedFile)">copy link</div>
        <div class="overlay__item file__delete" v-on:click="deleteFile(selectedFile)">delete</div>
      </div>
    </Overlay>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, Ref } from "vue"
import Overlay from "@/components/ui/Overlay.vue"
import FileAdapter from "@/common/adapters/Main/FileAdapter"
import Dates from "@/common/tools/Dates"
import Sizes from "@/common/tools/Sizes.js"
import Extensions from "@/common/tools/Extensions"
import { IMeta, iMetaDefault } from '@/types/global'
import { IFile, IFileDefault } from '@/types/file'

const filesLoaded: Ref<boolean> = ref(false)
const isToolsOverlayActive: Ref<boolean> = ref(false)
const isSortOverlayActive: Ref<boolean> = ref(false)
const sortBy: Ref<string> = ref('created')
const sortFirst: Ref<string> = ref('newest')
const files: Ref<Array<IFile>> = ref([IFileDefault])
const meta: Ref<IMeta> = ref(iMetaDefault)
const selectedFile: Ref<IFile> = ref(IFileDefault)
const perPage: Ref<string> = ref('')
const next: Ref<string> = ref('')
// file input for upload
const input= ref(null)


onMounted(() => {
  getFiles('')
})

function getFiles(cursor: string) {
  filesLoaded.value = false
  FileAdapter.getFiles(cursor, undefined)
    .then(result => {
      files.value = result.data
      meta.value = result.meta
      next.value = meta.value.next
      filesLoaded.value = true
    })
}


async function refresh() {
  let isTrueFiles = filesLoaded.value && files.value.length < 1
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
  const isDelete = confirm('Удалить файл?')
  if (isDelete) {
    FileAdapter.delete(file.id).then(() => {
      deleteFileFromArray(file)
      isToolsOverlayActive.value = false
    })
  }
}

async function onFileInputChange(event) {
  const files = event.target.files
  if (files.length < 1) {
    return 0
  }
  await FileAdapter.upload(files)
  await getFiles('')
}

function onUploadClick() {
  if (!input.value) {
    return
  }
  input.value.value = ''
  input.value.click()
}

async function copyLink(file: IFile) {
  const url = `${import.meta.env.VITE_UPLOADS_URL}/${file.path}`
  navigator.clipboard.writeText(url)
    .then(() => {
      isToolsOverlayActive.value = false
      window.$elvenNotify.add('Link copied to clipboard.')
    })
}

function deleteFileFromArray(file: IFile) {
  const index = files.value.indexOf(file)
  files.value.splice(index, 1)
  return true
}

function select(file) {
  isToolsOverlayActive.value = true
  selectedFile.value = file
}

function convertDateWrap(date) {
  return Dates.convert(date)
}

function convertSizeWrap(size) {
  return Sizes.convert(size)
}

function readableExtensionWrap(extension) {
  return Extensions.getReadable(extension)
}

function convertPreviewPath(path) {
  return `${import.meta.env.VITE_UPLOADS_URL}/${path}`
}

function playAudio(file) {
  const converted = convertPreviewPath(file.path)
  window.$elvenPlayer.play(converted)
}
</script>

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

.overlay__file-tools {
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