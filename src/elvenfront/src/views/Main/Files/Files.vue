<template>
  <div class="files-main">
    <div class="files-tools">
      <div class="files-tools-upload cursor-pointer">
        <div class="files-tools-upload-butt" v-on:click="onUploadClick()">Загрузить</div>
        <input
          id="input-upload"
          type="file"
          multiple
          style="display: none"
          @input="onFileInputChange($event)"
        />
      </div>
    </div>
    <div class="files-list" v-if="filesLoaded && files.length > 0">
      <div class="file" v-for="file in files" :key="file.id" v-on:click="select(file)">
        <div class="file-meta">
          <div class="file-item file-loaded-at">{{ convertDateWrap(file.created_at) }}</div>
        </div>
        <div class="file-main">
          <div
            class="file-item file-preview"
            v-on:click.stop
            v-if="readableExtensionWrap(file.extension) === 'IMAGE'"
          >
            <img
              :src="convertPreviewPath(file.path)"
              v-if="readableExtensionWrap(file.extension) === 'IMAGE'"
            />
          </div>
          <div
            class="file-item file-preview"
            v-on:click.stop
            v-if="readableExtensionWrap(file.extension) === 'VIDEO'"
          >
            <video controls :src="convertPreviewPath(file.path)"></video>
          </div>

          <div class="file-item file-name">{{ file.original_name }}</div>
          <div class="file-item file-size">{{ convertSizeWrap(file.size) }}</div>
        </div>
      </div>
    </div>

    <div class="files-404" v-if="filesLoaded && files.length < 1">
      <div class="files-404-1">Нет файлов :(</div>
    </div>

    <UIOverlay v-bind:active="isToolsOverlayActive" v-on:deactivated="isToolsOverlayActive = false">
      <div class="overlay-file-tools">
        <div
          class="ov-item file-play"
          v-if="readableExtensionWrap(selectedFile.extension) === 'AUDIO'"
          v-on:click="playAudio(selectedFile)"
        >Воспроизвести</div>
        <div class="ov-item file-copy-link" v-on:click="copyLink(selectedFile)">Скопировать ссылку</div>
        <div class="ov-item file-delete" v-on:click="deleteFile(selectedFile)">Удалить</div>
      </div>
    </UIOverlay>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref, Ref} from "vue"
import UIOverlay from "@/components/_UI/UIOverlay.vue"
import FileAdapter from "@/common/adapters/Main/FileAdapter"
import Dates from "@/common/tools/Dates"
import Sizes from "@/common/tools/Sizes.js"
import Extensions from "@/common/tools/Extensions"
import { IMeta, iMetaDefault } from '@/types/response'
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
  const inputUpload = document.getElementById('input-upload') as HTMLInputElement
  if (!inputUpload) {
    return
  }
  inputUpload.value = ''
  inputUpload.click()
}

async function copyLink(file: IFile) {
  const url = `${import.meta.env.VITE_UPLOADS_URL}/${file.path}`
  navigator.clipboard.writeText(url)
    .then(() => {
      isToolsOverlayActive.value = false
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
.files-main {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.files-tools {
  text-decoration: underline;
  background-color: var(--color-level-1);
  font-size: 1rem;
  width: 100%;
  height: 54px;
  border-radius: var(--border-radius);
  display: flex;
  flex-direction: row;
  align-items: center;
}

.files-tools-upload {
  margin-left: 12px;
}

.files-list {
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

.file-item {
  font-size: 0.9rem;
  line-height: 1.5rem;
  margin-top: 8px;
  margin-left: 12px;
  margin-right: 12px;
}

.file-meta {
  display: flex;
  flex-direction: row;
  color: var(--color-text-inactive);
}

.file-main {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-name {
  font-size: 1.1rem;
  line-height: 2rem;
  letter-spacing: 0.0099rem;
}

.file-size {
  font-size: 1rem;
}

.files-404 {
  background-color: var(--color-level-1);
  height: 240px;
  border-radius: var(--border-radius);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 24px;
}

.overlay-file-tools {
  width: 100%;
}

.ov-item {
  height: 64px;
  width: 100%;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ov-item:hover {
  background-color: var(--color-hover);
}
</style>