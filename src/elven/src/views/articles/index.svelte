<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { querystring } from "svelte-spa-router";
  // ui
  import Pagination from "@/components/pagination.svelte";
  // article
  import type { Article, Params } from "@/types/articles";
  import { Show, By, Start } from "@/types/articles";
  import type { Meta } from "@/types";
  import CArticle from "../../components/article.svelte";
  import Utils from "@/tools";
  import NetworkArticle from "@/network/network_article";
  import ArticleActions from "@/components/article_actions.svelte";
  import ArticlesToolbars from "@/components/articles_toolbars.svelte";

  /** articles loaded? */
  let loaded = false;

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

  /** is article selected? */
  let isSelected = false;

  /** selected article */
  let selected: {
    counter: number | null;
    article: Article | null;
    mouseEvent: MouseEvent;
  } = { counter: null, article: null, mouseEvent: null };

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

  function onDeleted() {
    isSelected = false;
    deleteFromArray(selected.counter);
  }

  /** when page changed */
  function onPageChanged(page: number) {
    requestParams.page = page;
    getAll();
  }

  /** on request param changed */
  async function onParamChanged() {
    await getAll();
  }

  /** select article */
  function select(article: Article, mouseEvent: MouseEvent, counter: number) {
    selected.counter = counter;
    selected.mouseEvent = mouseEvent;
    selected.article = articles[counter];
    isSelected = true;
  }

  /** get all articles */
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

  /** delete article from articles array and refresh articles */
  async function deleteFromArray(counter: number) {
    delete articles[counter];
    articles = articles;
    await refresh();
  }

  /** refresh articles */
  async function refresh() {
    const getData = async () => {
      await getAll()
      return articles
    }
    const setPage = (val: number) => (requestParams.page = val);
    await Utils.refresh(requestParams.page, setPage, getData);
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
  <ArticlesToolbars
    bind:params={requestParams}
    on:paramChanged={async () => await onParamChanged()}
  />

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
</style>
