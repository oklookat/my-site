<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { params, push } from "svelte-spa-router";
  // ui
  import Overlay from "@/ui/overlay.svelte";
  import Pagination from "@/ui/pagination.svelte";
  import Toolbar from "@/ui/toolbar.svelte";
  import ToolbarBig from "@/ui/toolbarBig.svelte";
  // article
  import type { Article, Params } from "./../types";
  import { Show, By, Start } from "./../types";
  import type { Meta } from "@/types";
  import CArticle from "./article.svelte";
  import Utils from "@/tools/utils";
  import ArticleAdapter from "../adapter";
  import Selector from "../../categories/comps/selector.svelte";
  import type { Category } from "../../categories/types";

  /** articles loaded? */
  let loaded = false;
  /** article selected / overlay opened? */
  let toolsOverlay = false;
  /** selected article */
  let selected: {
    counter: number | null;
    article: Article | null;
  } = { counter: null, article: null };
  /** response articles */
  let articles: Record<number, Article> = {};
  /** response information */
  let meta: Meta;
  /** request params */
  let requestParams: Params = {
    page: 1,
    show: Show.published,
    by: By.updated,
    start: Start.newest,
    preview: true,
    category_name: null,
    without_category: false,
  };
  /** custom categories for filter */
  const customCategories: Record<number, Category> = {
    "-1": {
      id: "-1",
      name: "All",
    },
  };

  onMount(async () => {
    await getAll(undefined);
  });

  onDestroy(() => {});

  /** get all articles.
   * @param p request params
   */
  async function getAll(p: Params = requestParams) {
    // trigger reactivity because baseParams may have been changed by history
    requestParams = requestParams;
    if (p.page < 1) {
      p.page = 1;
      requestParams.page = 1;
    }
    loaded = false;
    try {
      const result = await ArticleAdapter.getAll(p);
      articles = result.data;
      meta = result.meta;
      loaded = true;
    } catch (err) {}
  }

  /** edit article */
  function edit(counter: number) {
    const articleID = getIDByCounter(counter);
    push(`/articles/create/${articleID}`);
  }

  /** select article */
  function select(counter: number) {
    toolsOverlay = true;
    selected.counter = counter;
    selected.article = articles[counter];
  }

  /** publish article */
  async function publish(counter: number) {
    const converted = getIDByCounter(counter);
    toolsOverlay = false;
    try {
      await ArticleAdapter.publish(converted);
      await deleteFromArray(counter);
    } catch (err) {}
  }

  /** unpublish article */
  async function unpublish(counter: number) {
    const converted = getIDByCounter(counter);
    toolsOverlay = false;
    try {
      await ArticleAdapter.unpublish(converted);
      await deleteFromArray(counter);
    } catch (err) {}
  }

  /** delete article */
  async function deleteArticle(counter: number) {
    const isDelete = await window.$choose.confirm("delete article");
    if (!isDelete) {
      return;
    }
    toolsOverlay = false;
    try {
      const converted = getIDByCounter(counter);
      await ArticleAdapter.delete(converted);
      await deleteFromArray(counter);
    } catch (err) {}
  }

  /** delete article from articles array and refresh articles */
  async function deleteFromArray(counter: number) {
    delete articles[counter];
    articles = articles;
    await refresh();
  }

  /** set 'by' param and get articles */
  function setBy(by: By = By.published) {
    requestParams.by = by;
    requestParams.page = 1;
    getAll();
  }

  /** set 'start' param and get articles */
  function setStart(start: Start = Start.newest) {
    requestParams.start = start;
    requestParams.page = 1;
    getAll();
  }

  /** set 'show' param and get articles */
  function setShow(show: Show) {
    requestParams.show = show;
    requestParams.page = 1;
    getAll();
  }

  /** when page changed */
  function onPageChanged(page: number) {
    requestParams.page = page;
    getAll();
  }

  /** get article id by articles counter id */
  function getIDByCounter(counter: number): string {
    return articles[counter].id;
  }

  /** refresh articles */
  async function refresh() {
    while (true) {
      const articlesLength = Utils.getObjectLength(articles);
      const noArticles = loaded && articlesLength < 1;
      if (!noArticles) {
        break;
      }
      requestParams.page--;
      try {
        await getAll();
      } catch (err) {
        break;
      }
      if (requestParams.page < 2) {
        break;
      }
    }
  }

  /** sort by category */
  function onCategoryChanged(cat: Category | null) {
    requestParams.page = 1;
    // no categories
    requestParams.without_category = cat === null;
    let catName = null;
    const notAllCategories = cat && cat["name"] && cat.id !== "-1"
    if (notAllCategories) {
      catName = cat.name;
    }
    requestParams.category_name = catName;
    getAll();
  }
</script>

<div class="articles base__container">
  <ToolbarBig>
    <a href="#/articles/create">new</a>
    <a href="#/articles/cats">categories</a>
  </ToolbarBig>

  <Toolbar>
    <Selector
      {customCategories}
      selectedID={"-1"}
      on:changed={(e) => onCategoryChanged(e.detail)}
    />
  </Toolbar>

  <Toolbar>
    <div class="articles__item">
      {#if requestParams.show === Show.published}
        <div class="pointer" on:click={() => setShow(Show.drafts)}>
          published
        </div>
      {/if}
      {#if requestParams.show === Show.drafts}
        <div class="pointer" on:click={() => setShow(Show.published)}>
          drafts
        </div>
      {/if}
    </div>
    <div class="articles__item">
      {#if requestParams.start === Start.newest}
        <div class="pointer" on:click={() => setStart(Start.oldest)}>
          newest
        </div>
      {/if}
      {#if requestParams.start === Start.oldest}
        <div class="pointer" on:click={() => setStart(Start.newest)}>
          oldest
        </div>
      {/if}
    </div>
    <div class="articles__item">
      {#if requestParams.by === By.updated}
        <div class="pointer" on:click={() => setBy(By.published)}>
          by updated date
        </div>
      {/if}
      {#if requestParams.by === By.published}
        <div class="pointer" on:click={() => setBy(By.created)}>
          by published date
        </div>
      {/if}
      {#if requestParams.by === By.created}
        <div class="pointer" on:click={() => setBy(By.updated)}>
          by created date
        </div>
      {/if}
    </div>
  </Toolbar>

  {#if articles}
    <div class="articles__list">
      {#each Object.entries(articles) as [id, article]}
        <CArticle {article} on:selected={(e) => select(parseInt(id, 10))} />
      {/each}
    </div>
  {/if}

  {#if meta}
    <Pagination
      total={meta.total_pages}
      current={meta.current_page}
      on:changed={(e) => onPageChanged(e.detail)}
    />
  {/if}

  {#if toolsOverlay}
    <Overlay onClose={() => (toolsOverlay = false)}>
      <div class="overlay">
        {#if selected && selected.article.is_published}
          <div
            class="overlay__item"
            on:click={() => unpublish(selected.counter)}
          >
            unpublish
          </div>
        {:else}
          <div class="overlay__item" on:click={() => publish(selected.counter)}>
            publish
          </div>
        {/if}
        <div class="overlay__item" on:click={() => edit(selected.counter)}>
          edit
        </div>
        <div
          class="overlay__item"
          on:click={() => deleteArticle(selected.counter)}
        >
          delete
        </div>
      </div>
    </Overlay>
  {/if}
</div>

<style lang="scss">
  .articles {
    &__list {
      height: fit-content;
      width: 100%;
      display: flex;
      flex-direction: column;
      min-height: 42px;
      gap: 12px;
    }
  }

  .overlay {
    width: 100%;
    display: flex;
    flex-direction: column;
    &__item {
      height: 64px;
      width: 100%;
      font-size: 1rem;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    &__item:hover {
      background-color: var(--color-hover);
    }
  }
</style>
