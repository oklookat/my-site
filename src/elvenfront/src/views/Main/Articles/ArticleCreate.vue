<template>
  <div class="container">
    <Header></Header>
    <div class="content">
      <textarea id="article-title"
                placeholder="Если коротко..."
                rows="1" maxlength="124"
                v-model="article.title"
                @input="autoSave">
      </textarea>
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
import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter"
import UIOverlay from "@/components/_UI/UIOverlay"
import TextareaResizer from "@/common/tools/TextareaResizer"
import EditorJS from '@editorjs/editorjs'
import Head from '@editorjs/header'
import ImageTool from '@editorjs/image'

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

      // SERVICE START //
      textareaResizer: undefined,
      // SERVICE END //
    }
  },
  async mounted() {
    const route = useRoute()
    const id = route.params.id
    await this.initEditor()
    if (id) {
      this.article.id = id
      const isSuccess = await this.initEditArticle()
      if (isSuccess) {
        await this.setEditorData()
      } else {
        this.isEditorInitialized = true
      }
    } else {
      this.isEditorInitialized = true
    }
    this.textareaResizer = new TextareaResizer('article-title')
    this.textareaResizer.start()
  },
  async unmounted() {
    await window.editor.destroy()
    this.textareaResizer.destroy()
  },
  methods: {
    async initEditor() {
      window.editor = new EditorJS({
        holder: 'editor',
        tools: {
          paragraph: {
            config: {
              placeholder: 'Если развернуто...'
            }
          },
          header: {
            class: Head,
            inlineToolbar: true,
            config: {
              placeholder: 'Заголовок',
              levels: [2, 3, 4],
              defaultLevel: 3
            }
          },
          image: {
            class: ImageTool,
            config: {
              endpoints: {
                byFile: 'http://localhost:8008/uploadFile', // Your backend file uploader endpoint
                byUrl: 'http://localhost:8008/fetchUrl', // Your endpoint that provides uploading by Url
              }
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
          .catch(async (error) => {
            if (error === 404) {
              console.log('Запись не найдена')
              this.errorOverlayContent = 'Вы хотите отредактировать запись, которой не существует. Мы перенаправили вас на создание новой записи.'
            } else {
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
  gap: 32px;
}

.content {
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

#article-title {
  text-indent: 0;
  border: none;
  background-color: transparent;
  font-size: 1.6rem;
  font-weight: bold;
}

.editor-container {
  height: 100%;
}

#editor {
  height: 100%;
}

.error-ok-button {
  border-radius: var(--border-radius);
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