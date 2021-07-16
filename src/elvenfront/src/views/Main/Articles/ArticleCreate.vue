<template>
  <div class="container">
    <Header title="новая запись"></Header>
    <div class="content">
      <input id="article-title" type="text" v-model="article.title" @input="autoSave"/>
      <div class="editor-container">
        <div id="editor" @input="autoSave">
        </div>
      </div>
    </div>


    <UIOverlay v-bind:active="isErrorOverlayActive" v-on:deactivated="isErrorOverlayActive = false">
      {{ errorOverlayContent }}
      <div class="error-ok-button" v-on:click="isErrorOverlayActive = false">Ок</div>
    </UIOverlay>

  </div>
</template>

<script>
import {defineComponent} from "vue"
import {useRoute} from 'vue-router'
import Header from "@/components/Header/Header"
import EditorJS from '@editorjs/editorjs'
import Head from '@editorjs/header'
import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter"
import UIOverlay from "@/components/_UI/UIOverlay";

export default defineComponent({
  name: 'ArticleCreate',
  components: {UIOverlay, Header},
  data() {
    return {
      article: {
        id: null,
        title: '',
        content: undefined,
      },
      timeoutID: undefined,
      editorTimeoutID: undefined,
      editorLoadAttempts: 0,
      isEditorInitialized: false,
      isErrorOverlayActive: false,
      errorOverlayContent: '',
    }
  },
  async mounted() {
    const route = useRoute()
    const id = route.params.id
    await this.initEditor()
    if (id) {
      this.article.id = id
      const isSuccess = await this.initEditArticle()
      if(isSuccess){
        await this.setEditorData()
      } else {
        this.isEditorInitialized = true
      }
    } else {
      this.isEditorInitialized = true
    }
  },
  async unmounted() {
    await window.editor.destroy()
  },
  methods: {
    async initEditor() {
      window.editor = new EditorJS({
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
      })
    },
    async autoSave() {
      if (!this.isEditorInitialized) {
        throw Error('editor not initialized.')
      }
      if (this.timeoutID) {
        clearTimeout(this.timeoutID)
      }
      this.timeoutID = setTimeout(async () => {
        let saveAllowed = false
        await window.editor.save().then((outputData) => {
          if (outputData.blocks.length >= 1) {
            saveAllowed = true
            this.article.content = outputData
          }
        })
        if (!saveAllowed) {
          return 0
        }
        if (!this.article.id) {
          await ArticleAdapter.createArticle(this.article)
              .then(article => {
                this.article.id = article.id
              })
        } else {
          await ArticleAdapter.saveArticle(this.article)
        }
      }, 1000)
    },
    // ARTICLE EDITING FUNCTIONS START //
    async initEditArticle() {
      return await ArticleAdapter.getArticle(this.article.id)
          .then((article) => {
            this.article = article
            return Promise.resolve(true)
          })
          .catch(async (error) =>{
            if(error === 404){
              console.log('Запись не найдена')
              this.errorOverlayContent = 'Вы хотите отредактировать запись, которой не существует. Мы перенаправили вас на создание новой записи.'
            } else{
              this.errorOverlayContent = `Произошла странная ошибка. Ошибка: ${error}`
            }
            this.article.id = undefined
            await this.$router.push({name: 'ArticleCreate'})
            this.isErrorOverlayActive = true
            return Promise.resolve(false)
          })
    },
    async setEditorData() {
      this.editorTimeoutID = setInterval(async () => {
        if (this.editorLoadAttempts > 20) {
          this.isEditorInitialized = false
          clearInterval(this.editorTimeoutID)
          throw Error('failed while trying load data to editor. Is internet down? Server problems? Or article content have bad format, like bad/not JSON parsed.')
        }
        try {
          await window.editor.render(JSON.parse(this.article.content))
          this.isEditorInitialized = true
          clearInterval(this.editorTimeoutID)
        } catch (err) {
          console.info('Trying to load data in editor...')
          this.editorLoadAttempts++
        }
      }, 500)
    },
    // ARTICLE EDITING FUNCTIONS END //
  },
})
</script>

<style scoped>
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

.editor-container{
  height: 100%;
}


.error-ok-button{
  border-radius: 6px;
  background-color: var(--color-text);
  width: 25%;
  height: 10%;
  color: var(--color-body);
  cursor: pointer;
  font-weight: bold;
  font-size: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>