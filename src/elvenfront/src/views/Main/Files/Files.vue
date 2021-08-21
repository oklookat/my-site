<template>
  <div class="files-main">
    <div class="files-tools">
      <div class="files-tools-upload cursor-pointer">
        <div class="files-tools-upload-butt" v-on:click="onUploadClick">Загрузить</div>
        <input id="input-upload" type="file" multiple style="display: none" @input="onFileInputChange"/>
      </div>
    </div>

    <div class="files-list" v-if="isFilesLoaded && files.length > 0">
      <div class="file" v-for="file in files" :key="file.id" v-on:click="selectFile(file)">
        <div class="file-meta">
          <div class="file-item file-loaded-at">
            {{ convertDateWrap(file.created_at) }}
          </div>
        </div>
        <div class="file-main">
          <div class="file-item file-preview" v-on:click.stop>
            <img :src="convertPreviewPath(file.path)" v-if="readableExtensionWrap(file.extension) === 'IMAGE'">
            <video controls :src="convertPreviewPath(file.path)"
                   v-if="readableExtensionWrap(file.extension) === 'VIDEO'"></video>
          </div>
          <div class="file-item file-name">{{ file.original_name }}</div>
          <div class="file-item file-size">{{ convertSizeWrap(file.size) }}</div>
        </div>
      </div>
    </div>

    <div class="files-404" v-if="isFilesLoaded && files.length < 1">
      <div class="files-404-1">Нет файлов :(</div>
    </div>

    <UIPagination
        :total-pages="totalPages"
        :current-page="currentPage"
        v-on:page-changed="getFiles($event)">
    </UIPagination>


    <UIOverlay v-bind:active="isToolsOverlayActive" v-on:deactivated="isToolsOverlayActive = false">
      <div class="overlay-file-tools">
        <div class="ov-item file-play"
             v-if="readableExtensionWrap(selectedFile.extension) === 'AUDIO'"
             v-on:click="playAudio(selectedFile)">Воспроизвести</div>
        <div class="ov-item file-copy-link" v-on:click="copyLink(selectedFile)">Скопировать ссылку</div>
        <div class="ov-item file-delete" v-on:click="deleteFile(selectedFile)">Удалить</div>
      </div>
    </UIOverlay>


  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import FileAdapter from "@/common/adapters/Main/FileAdapter"
import Dates from "@/common/tools/Dates"
import Sizes from "@/common/tools/Sizes.js"
import UIOverlay from "@/components/_UI/UIOverlay.vue"
import Extensions from "@/common/tools/Extensions"
import UIPagination from "@/components/_UI/UIPagination.vue"

export default defineComponent({
  name: 'Files',
  components: {UIPagination, UIOverlay},
  data() {
    return {
      isFilesLoaded: false,
      isToolsOverlayActive: false,
      isSortOverlayActive: false,
      sortBy: 'created', // see backend docs for more
      sortFirst: 'newest',
      files: [],
      filesMeta: [],
      selectedFile: undefined,
      totalPages: 1,
      currentPage: 1,
    }
  },
  async mounted() {
    await this.getFiles()
  },
  methods: {
    // GET FUNCTIONS START //
    async getFiles(page = this.currentPage, sortFirst = this.sortFirst) {
      if (this.currentPage < 1) {
        this.currentPage = 1
        page = this.currentPage
      }
      this.isFilesLoaded = false
      await FileAdapter.getFiles(page, sortFirst)
          .then(result => {
            this.files = result.data
            this.filesMeta = result.meta
            this.perPage = this.filesMeta.per_page
            this.currentPage = this.filesMeta.current_page
            this.totalPages = Math.ceil(this.filesMeta.total / this.filesMeta.per_page)
            this.isFilesLoaded = true
          })
    },
    async refreshFiles() {
      let isTrueFiles = this.isFilesLoaded && this.files.length < 1
      if (isTrueFiles) { // no files in current page
        while (isTrueFiles) {
          // moving back until the pages ends or data appears
          this.currentPage--
          await this.getFiles()
          if (this.currentPage <= 1) {
            break
          }
          isTrueFiles = this.isFilesLoaded && this.files.length < 1
        }
      }
    },
    async deleteFile(file) {
      const isDelete = confirm('Удалить файл?')
      if (isDelete) {
        await FileAdapter.delete(file.id)
        this.deleteFileFromArray(file)
        this.isToolsOverlayActive = false
        await this.refreshFiles()
      }
    },
    // GET FUNCTIONS END //

    // UPLOAD FUNCTIONS START //
    async onFileInputChange(event) {
      const files = event.target.files
      if (files.length < 1) {
        return 0
      }
      await FileAdapter.upload(files)
      await this.getFiles()
    },
    onUploadClick() {
      const inputUpload = document.getElementById('input-upload')
      inputUpload.value = ''
      inputUpload.click()
    },
    // UPLOAD FUNCTIONS END //

    // SERVICE START //
    async copyLink(file) {
      const url = `${import.meta.env.VITE_UPLOADS_URL}/${file.path}`
      await navigator.clipboard.writeText(url)
          .then(() => {
            this.isToolsOverlayActive = false
          })
    },
    deleteFileFromArray(file) {
      const index = this.files.indexOf(file)
      this.files.splice(index, 1)
      return true
    },
    selectFile(file) {
      this.isToolsOverlayActive = true
      this.selectedFile = file
    },
    convertDateWrap(date) {
      return Dates.convert(date)
    },
    convertSizeWrap(size) {
      return Sizes.convert(size)
    },
    readableExtensionWrap(extension) {
      return Extensions.getReadable(extension)
    },
    convertPreviewPath(path) {
      return `${import.meta.env.VITE_UPLOADS_URL}/${path}`
    },
    playAudio(file){
      const converted = this.convertPreviewPath(file.path)
      this.$elvenPlayer.play(converted)
    }
    // SERVICE END //
  },
})
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
  min-height: 164px;
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
  font-size: 1rem;
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