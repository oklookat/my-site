<script lang="ts">
  import EditorJS from "@editorjs/editorjs";
  import Header from "@editorjs/header";
  import ImageTool from "@editorjs/image";
  import ArticleAdapter from "@/adapters/ArticleAdapter";
  import TextareaResizer from "@/tools/TextareaResizer";
  import type { TArticle } from "@/types/article";
  import { onDestroy, onMount } from "svelte";

  export let params: { id?: string } = { id: null };

  let editor: EditorJS;
  let autoSaveThrottle: ReturnType<typeof setTimeout> | null = null;
  let textareaResizer: TextareaResizer;
  // title element for textarea resize
  let articleTitleEL: HTMLElement | null;
  let article: TArticle = {
    title: "",
    content: null,
  };

  onMount(async () => {
    textareaResizer = new TextareaResizer(articleTitleEL);
    let id = params.id;
    // if edit mode
    if (id) {
      try {
        article = await get(id);
      } catch (err) {}
    }
    initEditor(article.content);
  });

  onDestroy(() => {
    editor.destroy();
    textareaResizer.destroy();
  });

  function initEditor(data?: any) {
    editor = new EditorJS({
      holder: "editor",
      tools: {
        paragraph: {
          config: {
            placeholder: `It's a long story.`,
          },
        },
        header: {
          class: Header,
          inlineToolbar: true,
          config: {
            placeholder: "Header",
            levels: [2],
            defaultLevel: 2,
          },
        },
        image: {
          class: ImageTool,
          config: {
            endpoints: {
              byFile: "http://localhost:8008/uploadFile", // Your backend file uploader endpoint
              byUrl: "http://localhost:8008/fetchUrl", // Your endpoint that provides uploading by Url
            },
          },
        },
      },
      data,
    });
  }

  async function get(id: string) {
    try {
      const result = await ArticleAdapter.get(id);
      return Promise.resolve(result);
    } catch (err) {
      return Promise.reject();
    }
  }

  async function autoSave() {
    if (autoSaveThrottle) {
      clearTimeout(autoSaveThrottle);
    }
    autoSaveThrottle = setTimeout(async () => {
      let empty = false;
      await editor.save().then((outputData) => {
        empty = outputData.blocks.length < 1;
        article.content = outputData;
      });
      // not save if empty
      if (empty) {
        return;
      }
      // if not saved before (new article)
      if (!article.id) {
        ArticleAdapter.create(article).then((newArticle) => {
          article.id = newArticle.id;
        });
        return;
      }
      ArticleAdapter.update(article);
    }, 1000);
  }
</script>

<div class="create__container">
  <div class="create">
    <textarea
      id="title"
      placeholder="Actually..."
      rows="1"
      maxlength="124"
      bind:value={article.title}
      bind:this={articleTitleEL}
      on:input={() => autoSave()}
    />

    <div class="editor__container">
      <div id="editor" on:input={() => autoSave()} />
    </div>
  </div>
</div>

<style lang="scss">
  textarea::placeholder {
    color: black;
  }

  .create {
    width: 85%;
    height: 95%;
    margin: auto;
    display: flex;
    flex-direction: column;
    gap: 24px;
    &__container {
      max-width: 732px;
      margin: auto;
      background-color: white;
      color: black;
      border-radius: 6px;
      padding-bottom: 24px;
      @media screen and (max-width: 1023px) {
        width: 95%;
        height: 95%;
      }
      @media screen and (max-width: 1919px) {
        width: 95%;
        max-width: 512px;
      }
    }
  }

  #title {
    margin-top: 12px;
    min-height: 38px;
    text-indent: 0;
    border: none;
    background-color: transparent;
    font-size: 1.6rem;
    font-weight: bold;
  }

  .editor__container {
    width: 100%;
  }

  #editor {
    height: 100%;
  }
</style>
