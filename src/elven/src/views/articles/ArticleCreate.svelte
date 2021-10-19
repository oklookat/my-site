<script lang="ts">
  import EditorJS from "@editorjs/editorjs"
  import Head from "@editorjs/header"
  import ImageTool from "@editorjs/image"
  import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter"
  import type TextareaResizer from "@/common/tools/TextareaResizer"
  import type { IArticle } from "@/types/article"
  import Overlay from "@/components/ui/Overlay.svelte"
  import { onDestroy, onMount } from "svelte"

  export let article: IArticle = {
    id: "",
    user_id: "",
    is_published: false,
    title: "",
    content: null,
    slug: "",
    published_at: "",
    updated_at: "",
  };
  export let timeoutID: ReturnType<typeof setTimeout> | null = null;
  export let editorTimeoutID: ReturnType<typeof setTimeout> | null = null;
  export let editorLoadAttempts = 0;
  export let isEditorInitialized = false;
  export let isErrorOverlayActive = false;
  export let errorOverlayContent = "";
  let textareaResizer: TextareaResizer;

  onMount(async () => {
    // const route = useRoute();
    // let id = null;
    // if (route.params.id) {
    //   id = route.params.id.toString();
    // }
    // await initEditor();
    // if (id) {
    //   article.id = id;
    //   const isSuccess = await initEditArticle();
    //   if (isSuccess) {
    //     await setEditorData();
    //   } else {
    //     isEditorInitialized = true;
    //   }
    // } else {
    //   isEditorInitialized = true;
    // }
    // textareaResizer = new TextareaResizer("article-title");
    // textareaResizer.start();
  });

  onDestroy(async () => {
    await window.editor.destroy();
    textareaResizer.destroy();
    if (editorTimeoutID) {
      clearTimeout(editorTimeoutID);
    }
  });

  function initEditor() {
    window.editor = new EditorJS({
      holder: "editor",
      tools: {
        paragraph: {
          config: {
            placeholder: `It's a long story.`,
          },
        },
        header: {
          class: Head,
          inlineToolbar: true,
          config: {
            placeholder: "Header",
            levels: [2, 3, 4],
            defaultLevel: 3,
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
      minHeight: 0,
      data: {},
    });
  }

  async function autoSave() {
    if (!isEditorInitialized) {
      throw Error("editor not initialized");
    }
    if (timeoutID && timeoutID) {
      clearTimeout(timeoutID);
    }
    timeoutID = setTimeout(async () => {
      let saveAllowed = false;
      await window.editor.save().then((outputData) => {
        if (outputData.blocks.length >= 1) {
          saveAllowed = true;
          article.content = outputData;
        }
      });
      if (!saveAllowed) {
        return;
      }
      if (article.id.length < 1) {
        await ArticleAdapter.createArticle(article).then((newArticle) => {
          article.id = newArticle.id;
        });
      } else {
        await ArticleAdapter.saveArticle(article);
      }
    }, 1000);
  }

  async function initEditArticle() {
    if (article.id.length < 1) {
      return Promise.resolve(false);
    }
    try {
      const result = await ArticleAdapter.getArticle(article.id);
      article = result;
      return Promise.resolve(true);
    } catch (err) {
      if (err === 404) {
        errorOverlayContent =
          "Article not found. You are redirected to create a new article.";
      } else {
        errorOverlayContent = `Strange error: ${err}`;
      }
      article.id = "";
      await router.push({ name: "ArticleCreate" });
      isErrorOverlayActive = true;
      return Promise.resolve(false);
    }
  }

  async function setEditorData() {
    editorTimeoutID = setInterval(async () => {
      if (editorLoadAttempts > 20) {
        isEditorInitialized = false;
        clearInterval(editorTimeoutID);
        throw Error(
          "failed while trying load data to editor. Is internet down? Server problems? Or article content have bad format, like bad/not JSON parsed."
        );
      }
      try {
        if (article.id.length < 1) {
          clearInterval(editorTimeoutID);
          return;
        }
        await window.editor.render(article.content);
        isEditorInitialized = true;
        clearInterval(editorTimeoutID);
      } catch (err) {
        console.info(`Trying to load data in editor... Last error: ${err}`);
        editorLoadAttempts++;
      }
    }, 500);
  }
</script>

<main>
  <div class="articles-create-container">
    <div class="articles-create-main">
      <textarea
        id="article-title"
        placeholder="Actually..."
        rows="1"
        maxlength="124"
        v-model="article.title"
        on:input={autoSave}
      />
      <div class="editor-container">
        <div id="editor" on:input={autoSave} />
      </div>

      <Overlay
        v-bind:active="isErrorOverlayActive"
        v-on:deactivated="isErrorOverlayActive = false"
      >
        {{ errorOverlayContent }}
        <div
          class="error-ok-button"
          on:click={() => (isErrorOverlayActive = false)}
        >
          Ок
        </div>
      </Overlay>
    </div>
  </div>
</main>

<style>
  .articles-create-container {
    max-width: 732px;
    margin: auto;
    background-color: var(--color-level-1);
    border-radius: 6px;
    padding-bottom: 24px;
  }

  @media screen and (max-width: 1919px) {
    .articles-create-container {
      width: 95%;
      max-width: 512px;
    }
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
