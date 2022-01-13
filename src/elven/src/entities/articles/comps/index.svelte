<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  // ui
  import Overlay from "@/ui/overlay.svelte";
  import Pagination from "@/ui/pagination.svelte";
  import Toolbar from "@/ui/toolbar.svelte";
  import ToolbarBig from "@/ui/toolbarBig.svelte";
  // article
  import type { Article, Params } from "../types";
  import { Show, By, Start } from "../types";
  import type { Meta } from "@/types";
  import ArticleAdapter from "../adapter";
  import CArticle from "./article.svelte";
  import { Route } from "@/tools/paths";
  import Validate from "@/entities/articles/validate";

  let loaded = false;
  let toolsOverlay = false;
  let selected: Article;
  let articles: Array<Article> = [];
  let meta: Meta;
  let baseParams: Params = {
    page: 1,
    show: Show.published,
    by: By.updated,
    start: Start.newest,
    preview: true,
  };

  // main
  onMount(async () => {
    window.onpopstate = () => {
      const searchParams = Route.getSearchParams();
      Validate.params(baseParams, searchParams);
      getAll(undefined, false);
    };
    const searchParams = Route.getSearchParams();
    Validate.params(baseParams, searchParams);
    Route.setHistoryParams(baseParams);
    getAll();
  });

  async function getAll(p: Params = baseParams, withHistory: boolean = true) {
    if (p.page < 1) {
      p.page = 1;
    }
    loaded = false;
    try {
      const result = await ArticleAdapter.getAll(baseParams);
      articles = result.data;
      meta = result.meta;
      loaded = true;
      if (withHistory) {
        Route.setHistoryParams(baseParams);
      }
    } catch (err) {}
  }

  function edit(id: string) {
    push(`/articles/create/${id}`);
  }

  function select(a: Article) {
    toolsOverlay = true;
    selected = a;
  }

  async function publish(id: string) {
    try {
      await ArticleAdapter.publish(id);
      deleteFromArray(id);
      toolsOverlay = false;
      await refresh();
    } catch (err) {}
  }

  async function toDrafts(id: string) {
    try {
      await ArticleAdapter.makeDraft(id);
      deleteFromArray(id);
      toolsOverlay = false;
      await refresh();
    } catch (err) {}
  }

  async function deleteArticle(id: string) {
    const isDelete = confirm("Delete article?");
    if (!isDelete) {
      return;
    }
    try {
      await ArticleAdapter.delete(id);
      deleteFromArray(id);
      toolsOverlay = false;
    } catch (err) {}
  }

  function deleteFromArray(id: string) {
    articles = articles.filter((a) => a.id !== id);
    refresh();
  }

  function setBy(by: By = By.published) {
    baseParams.by = by;
    baseParams.page = 1;
    getAll();
  }

  function setStart(start: Start = Start.newest) {
    baseParams.start = start;
    baseParams.page = 1;
    getAll();
  }

  function setShow(show: Show) {
    baseParams.show = show;
    baseParams.page = 1;
    getAll();
  }

  function onPageChanged(page: number) {
    baseParams.page = page;
    getAll();
  }

  async function refresh() {
    let noArticles =
      loaded && (articles.length < 1 || articles.length < meta.per_page);
    while (noArticles) {
      baseParams.page--;
      await getAll();
      if (baseParams.page <= 1) {
        break;
      }
    }
  }
</script>

<div class="articles">
  <ToolbarBig>
    <a href="#/articles/create">new</a>
  </ToolbarBig>

  <Toolbar>
    <div class="articles__show">
      {#if baseParams.show === Show.published}
        <div
          class="articles__item articles__show-published"
          on:click={() => setShow(Show.drafts)}
        >
          published
        </div>
      {/if}
      {#if baseParams.show === Show.drafts}
        <div
          class="articles__item articles__show-drafts"
          on:click={() => setShow(Show.published)}
        >
          drafts
        </div>
      {/if}
    </div>
    <div class="articles__sort-by-date">
      {#if baseParams.start === Start.newest}
        <div class="articles__item" on:click={() => setStart(Start.oldest)}>
          newest
        </div>
      {/if}
      {#if baseParams.start === Start.oldest}
        <div class="articles__item" on:click={() => setStart(Start.newest)}>
          oldest
        </div>
      {/if}
    </div>
    <div class="articles__sort-by">
      {#if baseParams.by === By.updated}
        <div class="articles__item" on:click={() => setBy(By.published)}>
          by updated date
        </div>
      {/if}
      {#if baseParams.by === By.published}
        <div class="articles__item" on:click={() => setBy(By.created)}>
          by published date
        </div>
      {/if}
      {#if baseParams.by === By.created}
        <div class="articles__item" on:click={() => setBy(By.updated)}>
          by created date
        </div>
      {/if}
    </div>
  </Toolbar>

  {#if articles && articles.length > 0}
    <div class="articles__list">
      {#each articles as article (article.id)}
        <CArticle {article} on:selected={(e) => select(e.detail)} />
      {/each}
    </div>
  {/if}

  {#if loaded && articles.length < 1}
    <div class="articles__404">
      <div>no articles :(</div>
    </div>
  {/if}

  {#if meta}
    <Pagination
      total={meta.total_pages}
      current={meta.current_page}
      on:changed={(e) => onPageChanged(e.detail)}
    />
  {/if}

  <Overlay
    bind:active={toolsOverlay}
    on:deactivated={() => (toolsOverlay = false)}
  >
    <div class="overlay">
      {#if selected && selected.is_published}
        <div class="overlay__item" on:click={() => toDrafts(selected.id)}>
          make a draft
        </div>
      {:else}
        <div class="overlay__item" on:click={() => publish(selected.id)}>
          publish
        </div>
      {/if}
      <div class="overlay__item" on:click={() => edit(selected.id)}>edit</div>
      <div class="overlay__item" on:click={() => deleteArticle(selected.id)}>
        delete
      </div>
    </div>
  </Overlay>
</div>

<style lang="scss">
  .articles {
    width: 95%;
    height: 100%;
    max-width: 512px;
    margin: auto;
    display: flex;
    flex-direction: column;
    gap: 14px;
    &__item {
      cursor: pointer;
    }
    &__404 {
      background-color: var(--color-level-1);
      height: 240px;
      border-radius: var(--border-radius);
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 24px;
    }
    &__list {
      height: 100%;
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
      transition: background-color 80ms ease-out;
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
