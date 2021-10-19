<template>
  <div class="files__container">
    <div class="files__tools">
      <div class="file__upload">
        <div class="file__upload-butt" v-on:click="onUploadClick()">upload</div>
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

    <div class="files__404" v-if="filesLoaded && files.length < 1">
      <div>no files :(</div>
    </div>

    <FilesList :files="files" @delete="deleteFile($event)"></FilesList>

    <Pagination :total-pages="totalPages" :current-page="currentPage" @changed="getFiles($event)"></Pagination>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, Ref } from "vue"
import FileAdapter from "@/common/adapters/Main/FileAdapter"
import { IMeta, iMetaDefault } from '@/types/global'
import { IFile } from '@/types/file'
import FilesList from "@/components/parts/FilesList.vue"
import Pagination from "@/components/ui/Pagination.vue"

const filesLoaded: Ref<boolean> = ref(false)
// service
const isSortOverlayActive: Ref<boolean> = ref(false)
const sortBy: Ref<string> = ref('created')
const sortFirst: Ref<string> = ref('newest')
// meta
const currentPage: Ref<number> = ref(1)
const totalPages: Ref<number> = ref(1)
// files
const files: Ref<Array<IFile>> = ref([])
const meta: Ref<IMeta> = ref(iMetaDefault)
const show: Ref<string> = ref('newest')
// file input for upload
const input = ref(null)


onMounted(() => {
  getFiles()
})

function getFiles(pageA = currentPage.value, showA = show.value) {
  currentPage.value = pageA
  show.value = showA
  filesLoaded.value = false
  FileAdapter.getFiles(pageA, showA)
    .then(result => {
      files.value = result.data
      meta.value = result.meta
      currentPage.value = meta.value.current_page
      totalPages.value = meta.value.total_pages
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
  const isDelete = confirm('delete file?')
  if (isDelete) {
    FileAdapter.delete(file.id).then(() => {
      deleteFileFromArray(file)
    })
  }
}

async function onFileInputChange(event) {
  const files = event.target.files
  if (files.length < 1) {
    return 0
  }
  await FileAdapter.upload(files)
  await getFiles()
}

function onUploadClick() {
  if (!input.value) {
    return
  }
  input.value.value = ''
  input.value.click()
}

function deleteFileFromArray(file: IFile) {
  const index = files.value.indexOf(file)
  files.value.splice(index, 1)
  return true
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