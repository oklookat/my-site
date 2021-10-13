<template>
  <div class="articles-create-container">
    <div class="articles-create-main">
      <textarea
        id="article-title"
        placeholder="Если коротко..."
        rows="1"
        maxlength="124"
        v-model="article.title"
        @input="autoSave"
      ></textarea>
      <div class="editor-container">
        <div id="editor" @input="autoSave"></div>
      </div>

      <UIOverlay
        v-bind:active="isErrorOverlayActive"
        v-on:deactivated="isErrorOverlayActive = false"
      >
        {{ errorOverlayContent }}
        <div class="error-ok-button" v-on:click="isErrorOverlayActive = false">Ок</div>
      </UIOverlay>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineComponent } from "vue"
import { useRoute, useRouter } from 'vue-router'
import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter"
import UIOverlay from "@/components/_UI/UIOverlay.vue"
import TextareaResizer from "@/common/tools/TextareaResizer"
import EditorJS from '@editorjs/editorjs'
import Head from '@editorjs/header'
import ImageTool from '@editorjs/image'
import { Ref, ref, reactive } from "@vue/reactivity"
import { onMounted, onUnmounted } from "@vue/runtime-core"
import { IArticle } from "@/types/article"

const router = useRouter()
const route = useRoute()
const article: Ref<IArticle> = ref({
  id: '',
  user_id: '',
  is_published: false,
  title: '',
  content: '',
  slug: '',
  published_at: '',
  updated_at: ''
})
const timeoutID: Ref<ReturnType<typeof setTimeout> | null> = ref(null)
const editorTimeoutID: Ref<ReturnType<typeof setTimeout> | null> = ref(null)
const editorLoadAttempts = ref(0)
const isEditorInitialized = ref(false)
const isErrorOverlayActive = ref(false)
const errorOverlayContent = ref('')
let textareaResizer: TextareaResizer

onMounted(async () => {
  const route = useRoute()
  let id = null
  if(route.params.id) {
    id = route.params.id.toString()
  }
  await initEditor()
  if (id) {
    article.value.id = id
    const isSuccess = await initEditArticle()
    if (isSuccess) {
      await setEditorData()
    } else {
      isEditorInitialized.value = true
    }
  } else {
    isEditorInitialized.value = true
  }
  textareaResizer = new TextareaResizer('article-title')
  textareaResizer.start()
})

onUnmounted(async () => {
  await window.editor.destroy()
  textareaResizer.destroy()
  if (editorTimeoutID.value) {
    clearTimeout(editorTimeoutID.value)
  }
})

function initEditor() {
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
}

async function autoSave() {
  if (!isEditorInitialized) {
    throw Error('editor not initialized.')
  }
  if (timeoutID && timeoutID.value) {
    clearTimeout(timeoutID.value)
  }
  timeoutID.value = setTimeout(async () => {
    let saveAllowed = false
    await window.editor.save().then((outputData) => {
      if (outputData.blocks.length >= 1) {
        saveAllowed = true
        article.value.content = outputData
      }
    })
    if (!saveAllowed) {
      return
    }
    if (article.value.id.length < 1) {
      await ArticleAdapter.createArticle(article.value)
        .then(newArticle => {
          article.value.id = newArticle.id
        })
    } else {
      await ArticleAdapter.saveArticle(article.value)
    }
  }, 1000)
}

async function initEditArticle() {
  if (article.value.id.length < 1) {
    return Promise.resolve(false)
  }
  try {
    const result = await ArticleAdapter.getArticle(article.value.id)
    article.value = result
    return Promise.resolve(true)
  } catch (err) {
    if (err === 404) {
      console.log('Запись не найдена')
      errorOverlayContent.value = 'Вы хотите отредактировать запись, которой не существует. Мы перенаправили вас на создание новой записи.'
    } else {
      errorOverlayContent.value = `Произошла странная ошибка. Ошибка: ${err}`
    }
    article.value.id = ''
    await router.push({ name: 'ArticleCreate' })
    isErrorOverlayActive.value = true
    return Promise.resolve(false)
  }
}

async function setEditorData() {
  editorTimeoutID.value = setInterval(async () => {
    if (editorLoadAttempts.value > 20) {
      isEditorInitialized.value = false
      clearInterval(editorTimeoutID.value)
      throw Error('failed while trying load data to editor. Is internet down? Server problems? Or article content have bad format, like bad/not JSON parsed.')
    }
    try {
      if (article.value.id.length < 1) {
        clearInterval(editorTimeoutID.value)
        return
      }
      await window.editor.render(article.value.content)
      isEditorInitialized.value = true
      clearInterval(editorTimeoutID.value)
    } catch (err) {
      console.info(`Trying to load data in editor... Last error: ${err}`)
      editorLoadAttempts.value++
    }
  }, 500)
}
</script>

<style scoped>
.articles-create-container {
  width: 100%;
  background-color: white;
  border-radius: 6px;
  padding-bottom: 24px;
}
.articles-create-main {
  width: 85%;
  height: 95%;
  margin: auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

#article-title {
  margin-top: 12px;
  min-height: 38px;
  text-indent: 0;
  border: none;
  background-color: transparent;
  font-size: 1.6rem;
  font-weight: bold;
}

textarea {
  text-align: center;
}

.editor-container {
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

@media screen and (max-width: 1023px) {
  .articles-create-main {
    width: 95%;
    height: 95%;
  }
}
</style>