<template>
  <div class="container">
    <Header></Header>
    <div class="content">
      <div class="files-tools">
        <div class="files-tools-upload cursor-pointer">
          <div class="files-tools-upload-butt" v-on:click="onUploadClick">Загрузить</div>
          <div class="the-underline"></div>
          <input id="input-upload" type="file" multiple style="display: none" @input="onFileInputChange"/>
        </div>
      </div>

      <div class="files-list" v-if="isFilesLoaded && files.length > 0">
        <div class="file" v-for="file in files" :key="file.id">
          <div class="file-meta">
            <div class="file-item file-loaded-at">
              {{ convertDateWrap(file.created_at) }}
            </div>
          </div>
          <div class="file-main">
            <div class="file-item file-preview" v-if="isToolsOverlayActive"></div>
            <div class="file-item file-name">{{ file.original_name }}</div>
            <div class="file-item file-size">{{ convertSizeWrap(file.size) }}</div>
          </div>
        </div>
      </div>

      <div class="files-404" v-if="isFilesLoaded && files.length < 1">
        <div class="files-404-1">Нет файлов :(</div>
      </div>

    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue"
import Header from "@/components/Header/Header"
import FileAdapter from "@/common/adapters/Main/FileAdapter"
import Dates from "@/common/tools/Dates"
import Sizes from "@/common/tools/Sizes.js"

export default defineComponent({
  name: 'Files',
  components: {Header},
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
    async getFiles(page = this.currentPage, sortBy = this.sortBy, sortFirst = this.sortFirst) {
      if (this.currentPage < 1) {
        this.currentPage = 1
        page = this.currentPage
      }
      this.isFilesLoaded = false
      await FileAdapter.getFiles(page, sortBy, sortFirst)
          .then(result => {
            this.files = result.data
            this.filesMeta = result.meta
            this.perPage = this.filesMeta.per_page
            this.currentPage = this.filesMeta.current_page
            this.totalPages = Math.ceil(this.filesMeta.total / this.filesMeta.per_page)
            this.isFilesLoaded = true
          })
    },
    // GET FUNCTIONS END //

    // UPLOAD FUNCTIONS START //
    async onFileInputChange(event) {
      const files = event.target.files
      if (files.length < 1) {
        return 0
      }
      await FileAdapter.upload(files)
    },
    onUploadClick() {
      const inputUpload = document.getElementById('input-upload')
      inputUpload.value = ''
      inputUpload.click()
    },
    // UPLOAD FUNCTIONS END //

    // SERVICE START //
    convertDateWrap(date) {
      return Dates.convert(date)
    },
    convertSizeWrap(size){
      return Sizes.convert(size)
    }
    // SERVICE END //
  },
})
</script>

<style scoped>
.content {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.files-tools {
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
  font-size: 1rem;
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

.file-size{
  font-size: 1rem;
}
</style>