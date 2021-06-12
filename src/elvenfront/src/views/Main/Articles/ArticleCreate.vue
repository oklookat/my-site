<template>
  <div class="container">
    <Header title="новая запись"></Header>
    <div class="content">
      <div id="editor"></div>
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue"
import Header from "@/components/Header/Header"
import EditorJS from '@editorjs/editorjs'
import Head from '@editorjs/header'

export default defineComponent({
  name: 'ArticleCreate',
  components: {Header},
  data(){
    return{
      editor: EditorJS
    }
  },
  mounted() {
    this.initEditor()
  },
  unmounted() {
    this.editor.destroy()
  },
  methods: {
    initEditor(){
      this.editor = new EditorJS({
        holder : 'editor',

        tools: {
          header: {
            class: Head,
            inlineToolbar : true,
            config: {
              placeholder: 'Превед медвед!',
              levels: [2, 3, 4],
              defaultLevel: 3
            }
          },
        },

        data: {}
      });
    }
  },
})
</script>

<style scoped lang="scss">
.content{
  display: block;
}
#editor{
  height: inherit;
  z-index: 1;
}
</style>