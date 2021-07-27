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
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue"
import Header from "@/components/Header/Header"
import FileAdapter from "@/common/adapters/Main/FileAdapter"

export default defineComponent({
  name: 'Files',
  components: {Header},
  data(){
    return{
    }
  },
  methods: {
    onUploadClick(){
      const inputUpload = document.getElementById('input-upload')
      inputUpload.value = ''
      inputUpload.click()
    },
    async onFileInputChange(event){
      const files = event.target.files
      if(files.length < 1){
        return 0
      }
      console.log(files)
      await FileAdapter.upload(files)
    },
  },
})
</script>

<style scoped>
.files-tools{
  background-color: var(--color-level-1);
  font-size: 1rem;
  width: 100%;
  height: 54px;
  border-radius: var(--border-radius);
  display: flex;
  flex-direction: row;
  align-items: center;
}

.files-tools-upload{
  margin-left: 12px;
}
</style>