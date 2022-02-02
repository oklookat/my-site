<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  // editor
  import jmarkd from "@oklookat/jmarkd";
  import "@oklookat/jmarkd/styles";
  import type { config } from "@oklookat/jmarkd";
  // utils
  import TextareaResizer from "@/tools/textareaResizer";
  // ui
  import Toolbar from "../../../../ui/toolbar.svelte";
  // article
  import type { Article } from "../types";
  import ArticleAdapter from "../adapter";
  import Select from "../../../../ui/select.svelte";
  import CategoryAdapter from "../../categories/adapter";
  import type { Category } from "../../categories/types";
  import Validate from "../validate";

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
  /** categories from server */
  let categories: Record<number, Category> = {};
  /** fetched categories for select element */
  let catsSelectData: {
    value: string;
    text: string;
  }[] = [];
  /** currently selected category value (counter) */
  let selectedCategoryValue: string | undefined;

  onMount(async () => {
    // check edit mode
    let id = params.id;
    if (id) {
      await getArticle(id);
    }
    // get available categories
    await getCategories();
    // manually add title before creating TextareaResizer, for correct height in start
    articleTitleEL.value = article.title;
    textareaResizer = new TextareaResizer(articleTitleEL);
    initEditor(article.content);
    // all loaded - set display
    createContainer.style.display = "flex";
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

  /** get categories for categories selector */
  async function getCategories() {
    // get categories
    try {
      const result = await CategoryAdapter.getAll();
      categories = result.data;
    } catch (err) {
      return Promise.reject();
    }
    makeCategoriesSelectable();
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

  /** format categories from server and put result in catsSelectData */
  function makeCategoriesSelectable() {
    catsSelectData = [];

    // create "no category" item
    const noCategory = { value: "", text: "No category" };
    catsSelectData[0] = noCategory;

    // if article without category - set "no category" as default
    const articleCatID = article.category_id;
    if (!articleCatID) {
      selectedCategoryValue = "";
    }

    // format categories for select element
    for (const [counter, _category] of Object.entries(categories)) {
      const option = {
        value: counter,
        text: _category.name,
      };
      // same categories?
      const sameCats = articleCatID === _category.id;
      if (sameCats) {
        selectedCategoryValue = counter;
      }
      catsSelectData.push(option);
    }

    // render
    catsSelectData = catsSelectData;
  }

  /** when category on select element changed */
  function onCategoryChanged(counter?: string) {
    // no counter = no category
    let newCatID = null;
    if (counter) {
      const cat = getCategoryByCounter(counter);
      if (!cat) {
        return;
      }
      newCatID = cat.id;
    }

    // category not changed?
    const notChanged = article.id === newCatID;
    if (notChanged) {
      return;
    }

    const oldCatID = article.category_id;
    article.category_id = newCatID;
    save().catch(() => {
      // revert changes
      article.category_id = oldCatID;
      makeCategoriesSelectable();
    });
  }

  /** get category by categories counter */
  function getCategoryByCounter(counter?: string | number): Category | null {
    let cat: Category | null = null;
    if (!counter) {
      return cat;
    }
    try {
      const isString = typeof counter === "string";
      const counterInt = isString ? parseInt(counter, 10) : counter;
      cat = categories[counterInt] || null;
    } catch (err) {}
    return cat;
  }
</script>

<div class="create" bind:this={createContainer}>
  <Toolbar>
    <Select
      bind:options={catsSelectData}
      bind:selected={selectedCategoryValue}
      on:selected={(e) => onCategoryChanged(e.detail)}
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
    // after data loaded - display = flex
    display: none;
    height: 100%;
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
    margin: auto;
  }

  .title {
    background-color: white;
    color: black;
    border: none;
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
