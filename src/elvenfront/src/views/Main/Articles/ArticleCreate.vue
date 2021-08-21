<template>
  <div class="articles-create-main">
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


    <UIOverlay v-bind:active="isErrorOverlayActive" v-on:deactivated="isErrorOverlayActive = false">
      {{ errorOverlayContent }}
      <div class="error-ok-button" v-on:click="isErrorOverlayActive = false">Ок</div>
    </UIOverlay>

  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import {useRoute} from 'vue-router'
import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter"
import UIOverlay from "@/components/_UI/UIOverlay.vue"
import TextareaResizer from "@/common/tools/TextareaResizer"
import EditorJS from '@editorjs/editorjs'
import Head from '@editorjs/header'
import ImageTool from '@editorjs/image'

export default defineComponent({
  name: 'ArticleCreate',
  components: {UIOverlay},
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
              captionPlaceholder: 'Описание',
              endpoints: {
                byFile: 'http://localhost:8008/uploadFile', // Your backend file uploader endpoint
                byUrl: 'http://localhost:8008/fetchUrl', // Your endpoint that provides uploading by Url
              },
            }
          },
        },
        i18n: {
          messages: {
            ui: {
              "blockTunes": {
                "toggler": {
                  "Click to tune": "Нажмите, чтобы настроить",
                  "or drag to move": "или перетащите"
                },
              },
              "inlineToolbar": {
                "converter": {
                  "Convert to": "Конвертировать в"
                }
              },
              "toolbar": {
                "toolbox": {
                  "Add": "Добавить"
                }
              }
            },
            toolNames: {
              "Text": "Параграф",
              "Heading": "Заголовок",
              "List": "Список",
              "Warning": "Примечание",
              "Checklist": "Чеклист",
              "Quote": "Цитата",
              "Code": "Код",
              "Delimiter": "Разделитель",
              "Raw HTML": "HTML-фрагмент",
              "Table": "Таблица",
              "Link": "Ссылка",
              "Marker": "Маркер",
              "Bold": "Полужирный",
              "Italic": "Курсив",
              "InlineCode": "Моноширинный",
            },
            tools: {
              "warning": {
                "Title": "Название",
                "Message": "Сообщение",
              },
              "link": {
                "Add a link": "Вставьте ссылку"
              },
              "stub": {
                'The block can not be displayed correctly.': 'Блок не может быть отображен'
              },
              image: {
                'Select an Image': 'Загрузить изображение',
                'Couldn’t upload image. Please try another.': 'Не удалось загрузить изображение.'
              }
            },
            blockTunes: {
              "image": {

              },
              "delete": {
                "Delete": "Удалить"
              },
              "moveUp": {
                "Move up": "Переместить вверх"
              },
              "moveDown": {
                "Move down": "Переместить вниз"
              }
            },
          }
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
.articles-create-main {
  width: 100%;
  height: 95%;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

#article-title {
  min-height: 38px;
  text-indent: 0;
  border: none;
  background-color: transparent;
  font-size: 1.6rem;
  font-weight: bold;
}

.editor-container {
  height: 100%;
  width: 100%;
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