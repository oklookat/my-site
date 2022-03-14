<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { querystring } from "svelte-spa-router";
  // ui
  import Pagination from "@/components/pagination.svelte";
  import Toolbar from "@/components/toolbar.svelte";
  import ToolbarBig from "@/components/toolbar_big.svelte";
  // article
  import type { Article, Params } from "@/types/articles";
  import { Show, By, Start } from "@/types/articles";
  import type { Meta } from "@/types";
  import CArticle from "../../components/article.svelte";
  import Utils from "@/tools";
  import NetworkArticle from "@/network/network_article";
  import CategoriesSelector from "@/components/categories_selector.svelte";
  import type { Category } from "@/types/articles/categories";
  import ArticleActions from "@/components/article_actions.svelte";

  /** articles loaded? */
  let loaded = false;

  /** is article selected? */
  let isSelected = false;

  /** selected article */
  let selected: {
    counter: number | null;
    article: Article | null;
    mouseEvent: MouseEvent;
  } = { counter: null, article: null, mouseEvent: null };

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

  /** custom categories for searching by category */
  const customCategories: Record<number, Category> = {
    "-1": {
      id: "-1",
      name: "All",
    },
  };

  //////////// PARAMS TEST
  let urlParams: URLSearchParams;
  const queryUnsub = querystring.subscribe((val) => {
    urlParams = new URLSearchParams(val);
  });

  onMount(async () => {
    await queryStringToRequestParams();
    await getAll(undefined);
  });

  //////////// PARAMS TEST
  // TODO: validate requestParams before send to server,
  // changed querystring if sort changed
  async function queryStringToRequestParams() {
    if (!urlParams) {
      return;
    }
    return new Promise((resolve) => {
      urlParams.forEach((key, name) => {
        console.log(name, key);
        if (!requestParams[name]) {
          return;
        }
        console.log("e");
        requestParams[name] = key;
      });
      resolve(true);
    });
  }

  // const newURL =
  //   window.location.protocol +
  //   "//" +
  //   window.location.host +
  //   window.location.pathname +
  //   "?myNewUrlQuery=1";

  // window.history.pushState({ path: newURL }, "", newURL);
  ///////////////////////

  onDestroy(() => {
    queryUnsub();
  });

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
      const result = await NetworkArticle.getAll(p);
      articles = result.data;
      meta = result.meta;
      loaded = true;
    } catch (err) {}
  }

  /** select article */
  function select(article: Article, mouseEvent: MouseEvent, counter: number) {
    selected.counter = counter;
    selected.mouseEvent = mouseEvent;
    selected.article = articles[counter];
    isSelected = true;
  }

  function onDeleted() {
    isSelected = false;
    deleteFromArray(selected.counter);
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
    const notAllCategories = cat && cat["name"] && cat.id !== "-1";
    if (notAllCategories) {
      catName = cat.name;
    }
    requestParams.category_name = catName;
    getAll();
  }
</script>

{#if isSelected}
  <ArticleActions
    article={selected.article}
    mouseEvent={selected.mouseEvent}
    onDisabled={() => (isSelected = false)}
    onDeleted={() => onDeleted()}
  />
{/if}

<div class="articles base__container">
  <div class="toolbars">
    <ToolbarBig>
      <a href="#/articles/create">new</a>
      <a href="#/articles/cats">categories</a>
    </ToolbarBig>

    <Toolbar>
      <CategoriesSelector
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
  </div>

  <div class="list">
    {#if articles}
      {#each Object.entries(articles) as [counter, article]}
        <CArticle
          {article}
          onSelected={(article, event) =>
            select(article, event, parseInt(counter, 10))}
        />
      {/each}
    {/if}
  </div>

  <div class="pages">
    {#if meta}
      <Pagination
        total={meta.total_pages}
        current={meta.current_page}
        on:changed={(e) => onPageChanged(e.detail)}
      />
    {/if}
  </div>
</div>

<style lang="scss">
  .toolbars {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
</style>
