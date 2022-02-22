<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  // editor
  import jmarkd from "@oklookat/jmarkd";
  import "@oklookat/jmarkd/styles";
  import type { config } from "@oklookat/jmarkd";
  // utils
  import Animation from "@/tools/animation";
  import TextareaResizer from "@/tools/textareaResizer";
  // ui
  import Toolbar from "../../../../ui/toolbar.svelte";
  // article
  import type { Article } from "../types";
  import ArticleAdapter from "../adapter";
  import Validate from "../validate";
  import Selector from "../../categories/comps/selector.svelte";
  import type { Category } from "../../categories/types";

  /** edit article with id (url params) */
  export let params: { id?: string };

  /** save all data */
  const save = saver();
  /** working on this article */
  let article: Article = {
    title: "",
    content: "",
  };
  /** main container */
  let createContainer: HTMLDivElement;
  /** title element */
  let articleTitleEL: HTMLTextAreaElement;
  /** title resizer */
  let textareaResizer: TextareaResizer;
  /** text editor */
  let editor: jmarkd;
  let editorEL: HTMLDivElement;

  onMount(async () => {
    // check edit mode
    let id = params.id;
    if (id) {
      await getArticle(id);
    }
    // manually add title before creating TextareaResizer, for correct height in start
    textareaResizer = new TextareaResizer(articleTitleEL, 54);
    articleTitleEL.value = article.title;
    initEditor(article.content);
    // all loaded - set opacity
    // (not display, because it brokes title resizing on init)
    await Animation.fadeIn(createContainer);
  });

  onDestroy(() => {
    editor.destroy();
    textareaResizer.destroy();
  });

  /** get article */
  async function getArticle(id: string) {
    try {
      const result = await ArticleAdapter.get(id);
      article = result;
    } catch (err) {
      return Promise.reject();
    }
    return Promise.resolve();
  }

  /** create new article */
  async function createArticle() {
    const notValid =
      article.id ||
      !Validate.title(article.title) ||
      !Validate.content(article.content);
    if (notValid) {
      return;
    }
    try {
      const newArticle = await ArticleAdapter.create(article);
      article.id = newArticle.id;
    } catch (err) {
      return Promise.reject(err);
    }
    return Promise.resolve();
  }

  /** update existing article */
  async function updateArticle() {
    const notValid =
      !article.id ||
      !Validate.title(article.title) ||
      !Validate.content(article.content);
    if (notValid) {
      return;
    }
    return await ArticleAdapter.update(article);
  }

  /** start text editor */
  function initEditor(data?: string) {
    // TODO: add sanitizer
    const config: config = {
      container: editorEL,
      placeholder: `It's a long story...`,
      input: data,
    };
    editor = new jmarkd(config);
  }

  /** create save func */
  function saver() {
    let throttle: NodeJS.Timeout;

    // save logic
    const saver = async () => {
      const outputData = editor.save();
      article.content = outputData;
      // if saved before (update)
      if (article.id) {
        return await updateArticle();
      }
      // if not saved before (new article)
      return await createArticle();
    };

    // save data
    return (): Promise<null> => {
      if (throttle) {
        clearTimeout(throttle);
      }
      return new Promise((resolve, reject) => {
        throttle = setTimeout(async () => {
          try {
            await saver();
            resolve(null);
          } catch (err) {
            reject(err);
          }
        }, 1000);
      });
    };
  }

  function onCategoryChanged(newCat: Category | null) {
    // category not changed?
    const newCatNotEmpty = newCat && newCat.id;
    const notChanged =
      (!article.category_id && !newCat) ||
      (newCatNotEmpty && article.category_id === newCat.id);
    if (notChanged) {
      return;
    }
    const oldCatID = article.category_id;
    if (newCatNotEmpty) {
      article.category_id = newCat.id;
    } else {
      // no category
      article.category_id = null;
    }
    save().catch(() => {
      // revert changes
      article.category_id = oldCatID;
    });
  }
</script>

<div class="create" bind:this={createContainer}>
  <Toolbar>
    <Selector
      bind:selectedID={article.category_id}
      on:changed={(e) => onCategoryChanged(e.detail)}
    />
  </Toolbar>

  <textarea
    class="title"
    placeholder="Actually..."
    rows="1"
    maxlength="124"
    bind:value={article.title}
    bind:this={articleTitleEL}
    on:input={() => save()}
  />

  <div class="editor" bind:this={editorEL} on:input={() => save()} />
</div>

<style lang="scss">
  .create {
    // after data loaded - opacity = 1
    opacity: 0;
    display: flex;
    height: fit-content;
    width: 95%;
    max-width: 424px;
    margin: auto;
    @media screen and (min-width: 1365px) {
      max-width: 724px;
    }
    flex-direction: column;
    align-items: center;
    gap: 24px;
  }

  .title,
  .editor {
    box-sizing: border-box;
    border: var(--color-border) 1px solid;
    margin: auto;
  }

  .title {
    background-color: white;
    color: black;
    font-size: 1.6rem;
    font-weight: bold;
    min-height: 54px;
    width: 100%;
    border-radius: 8px;
    padding: 12px;
  }

  .editor {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
  }
</style>
