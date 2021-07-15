<template>
  <div class="container">
    <Header title="новая запись"></Header>
    <div class="content">
      <input id="article-title" type="text" v-model="title" @input="saveData"/>
      <div id="editor" @input="saveData"></div>
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue"
import Header from "@/components/Header/Header"
import EditorJS from '@editorjs/editorjs'
import Head from '@editorjs/header'
import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter";

export default defineComponent({
  name: 'ArticleCreate',
  components: {Header},
  data() {
    return {
      editor: EditorJS,
      article: undefined,
      title: '',
    }
  },
  mounted() {
    this.initEditor()
  },
  unmounted() {
    this.editor.destroy()
  },
  methods: {
    initEditor() {
      this.editor = new EditorJS({
        holder: 'editor',

        tools: {
          header: {
            class: Head,
            inlineToolbar: true,
            config: {
              placeholder: 'Превед медвед!',
              levels: [2, 3, 4],
              defaultLevel: 3
            }
          },
        },
        minHeight: 0,
        data: {}
      });
    },
    async saveData() {
      let editorData
      let saveAllowed = false
      await this.editor.save().then((outputData) => {
        if(outputData.blocks.length >= 1){
          saveAllowed = true
          editorData = outputData
        }
      }).catch((error) => {
        console.log('Saving failed: ', error)
      })
      if(!saveAllowed){
        return 0
      }
      if (!this.article) {
        let article = {title: this.title, content: editorData}
        await ArticleAdapter.createArticle(article)
      } else {
        await ArticleAdapter.saveArticle(this.article)
      }
    }
  },
})
</script>

<style scoped lang="scss">
.container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.content {
  display: grid;
  height: calc(100% - var(--header-height));
  grid-template-columns: 1fr;
  grid-template-rows: 48px 1fr;
}

#editor {
  height: 100%;
  z-index: 1;
}

#article-title {
  width: 100%;
}
</style>