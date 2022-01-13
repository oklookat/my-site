<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  // editor
  import jmarkd from "@oklookat/jmarkd";
  import "@oklookat/jmarkd/styles";
  import type { config } from "@oklookat/jmarkd";
  // utils
  import TextareaResizer from "@/tools/textareaResizer";
  // article
  import type { Article } from "../types"
  import ArticleAdapter from "../adapter"

  export let params: { id?: string } = { id: null };

  let editor: jmarkd;
  let editorEL: HTMLDivElement;
  let autoSaveThrottle: ReturnType<typeof setTimeout> | null = null;
  let textareaResizer: TextareaResizer;
  // title element for textarea resize
  let articleTitleEL: HTMLTextAreaElement;
  let article: Article = {
    title: "",
    content: null,
  };

  onMount(async () => {
    let id = params.id;
    // if edit mode
    if (id) {
      try {
        article = await get(id);
      } catch (err) {}
    }
    // manually add title before creating TextareaResizer, for correct height in start
    articleTitleEL.value = article.title;
    textareaResizer = new TextareaResizer(articleTitleEL);
    initEditor(article.content);
  });

  onDestroy(() => {
    editor.destroy();
    textareaResizer.destroy();
  });

  function initEditor(data?: string) {
    // TODO: add sanitizer
    const config: config = {
      container: editorEL,
      placeholder: `It's a long story...`,
      input: data,
    };
    editor = new jmarkd(config);
  }

  async function get(id: string) {
    try {
      const result = await ArticleAdapter.get(id);
      return Promise.resolve(result);
    } catch (err) {
      return Promise.reject();
    }
  }

  function autoSave() {
    if (autoSaveThrottle) {
      clearTimeout(autoSaveThrottle);
    }
    autoSaveThrottle = setTimeout(() => {
      let empty = false;
      const outputData = editor.save();
      console.log(outputData);
      empty = outputData.length < 1;
      article.content = outputData;
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

<div class="create">
  <textarea
    class="title"
    placeholder="Actually..."
    rows="1"
    maxlength="124"
    bind:value={article.title}
    bind:this={articleTitleEL}
    on:input={() => autoSave()}
  />

  <div class="editor" bind:this={editorEL} on:input={() => autoSave()} />
</div>

<style lang="scss">
  .create {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 24px;
  }

  .title,
  .editor {
    margin: auto;
  }

  .title {
    background-color: white;
    color: black;
    border: none;
    font-size: 1.6rem;
    font-weight: bold;
    //
    min-height: 54px;
    width: 95%;
    max-width: 424px;
    border-radius: 8px;
    padding: 12px;
    @media screen and (min-width: 1365px) {
      max-width: 724px;
    }
  }

  .editor {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
  }
</style>
