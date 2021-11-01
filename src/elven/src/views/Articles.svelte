<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import type {
    TArticle,
    TBy,
    TParams,
    TShow,
    TStart,
  } from "@/types/ArticleTypes";
  import type { IMeta } from "@/types/GlobalTypes";
  import ArticleAdapter from "@/adapters/ArticleAdapter";
  import Overlay from "@/components/ui/Overlay.svelte";
  import Pagination from "@/components/ui/Pagination.svelte";
  import ArticlesList from "@/components/parts/ArticlesList.svelte";

  let loaded = false;
  let toolsOverlay = false;
  let selected: TArticle | null = null;
  let articles: Array<TArticle> = [];
  let meta: IMeta | null = null;
  let params: TParams = {
    page: 1,
    show: "published",
    by: "updated",
    start: "newest",
    preview: false,
  };

  // main
  onMount(async () => {
    await getAll();
  });

  async function getAll(p?: TParams) {
    if (!p) {
      p = params;
    } else {
      params = p;
    }
    loaded = false;
    ArticleAdapter.getAll(params).then((result) => {
      articles = result.data;
      meta = result.meta;
      loaded = true;
    });
  }

  async function edit(id: string) {
    await push(`/articles/create/${id}`);
  }

  function select(a: TArticle) {
    toolsOverlay = true;
    selected = a;
  }

  async function publish(id: string) {
    await ArticleAdapter.publish(id).then(() => {
      deleteFromArray(id);
      toolsOverlay = false;
    });
    refresh();
  }

  async function toDrafts(id: string) {
    await ArticleAdapter.makeDraft(id).then(() => {
      deleteFromArray(id);
      toolsOverlay = false;
    });
    refresh();
  }

  async function deleteArticle(id: string) {
    const isDelete = confirm("Delete article?");
    if (isDelete) {
      await ArticleAdapter.delete(id);
      deleteFromArray(id);
      toolsOverlay = false;
      refresh();
    }
  }

  function deleteFromArray(id: string) {
    articles = articles.filter((a) => a.id !== id);
    if (articles.length < meta.per_page) {
      getAll();
    }
  }

  async function setBy(by: TBy = "published") {
    params.by = by;
    params.page = 1;
    await getAll();
  }

  async function setStart(start: TStart = "newest") {
    params.start = start;
    params.page = 1;
    await getAll();
  }

  async function setShow(show: TShow) {
    params.show = show;
    params.page = 1;
    await getAll();
  }

  async function onPageChanged(page: number) {
    params.page = page;
    await getAll();
  }

  async function refresh() {
    // refresh is need when for ex. you deleted all articles on current page
    // and we need to check, is data on current page exists?
    // if page > 1 and no data, we moving back (currentPage--) and get new articles
    while (loaded && articles.length < 1) {
      // moving back until the pages ends or data appears
      params.page--;
      await getAll();
      if (params.page <= 1) {
        break;
      }
    }
  }
</script>

<div class="articles__container">
  <div class="articles__create">
    <a href="#/articles/create">new</a>
  </div>
  <div class="articles__toolbar">
    <div class="articles__show">
      {#if params.show === "published"}
        <div
          class="articles__item articles__show-published"
          on:click={() => setShow("drafts")}
        >
          show published
        </div>
      {/if}
      {#if params.show === "drafts"}
        <div
          class="articles__item articles__show-drafts"
          on:click={() => setShow("published")}
        >
          show drafts
        </div>
      {/if}
    </div>
    <div class="articles__sort-by-date">
      {#if params.start === "newest"}
        <div class="articles__item" on:click={() => setStart("oldest")}>
          newest
        </div>
      {/if}
      {#if params.start === "oldest"}
        <div class="articles__item" on:click={() => setStart("newest")}>
          oldest
        </div>
      {/if}
    </div>
    <div class="articles__sort-by">
      {#if params.by === "updated"}
        <div class="articles__item" on:click={() => setBy("published")}>
          by updated date
        </div>
      {/if}
      {#if params.by === "published"}
        <div
          class="articles__item"
          v-if="sortBy === 'published'"
          on:click={() => setBy("created")}
        >
          by published date
        </div>
      {/if}
      {#if params.by === "created"}
        <div class="articles__item" on:click={() => setBy("updated")}>
          by created date
        </div>
      {/if}
    </div>
  </div>

  {#if articles && articles.length > 0}
    <ArticlesList {articles} on:selected={(e) => select(e.detail)} />
  {/if}

  {#if loaded && articles.length < 1}
    <div class="articles__404">
      <div>no articles :(</div>
    </div>
  {/if}

  {#if meta}
    <Pagination
      totalPages={meta.total_pages}
      currentPage={meta.current_page}
      on:changed={(e) => onPageChanged(e.detail)}
    />
  {/if}

  <Overlay
    bind:active={toolsOverlay}
    on:deactivated={() => (toolsOverlay = false)}
  >
    <div class="overlay__article-manage" v-if="isToolsOverlayActive">
      {#if selected && selected.is_published}
        <div
          class="overlay__item make__draft"
          on:click={() => toDrafts(selected.id)}
        >
          make a draft
        </div>
      {:else}
        <div
          class="overlay__item publish"
          on:click={() => publish(selected.id)}
        >
          publish
        </div>
      {/if}
      <div class="overlay__item edit" on:click={() => edit(selected.id)}>
        edit
      </div>
      <div
        class="overlay__item delete"
        on:click={() => deleteArticle(selected.id)}
      >
        delete
      </div>
    </div>
  </Overlay>
</div>

<style>
  .articles__container {
    width: 95%;
    height: 100%;
    max-width: 512px;
    margin: auto;
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .articles__create {
    background-color: var(--color-level-1);
    height: 36px;
    width: 100%;
    display: flex;
    flex-direction: row;
  }

  .articles__create > a {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
  }

  .articles__create > a:hover {
    background-color: var(--color-hover);
  }

  .articles__toolbar {
    background-color: var(--color-level-1);
    padding-left: 12px;
    font-size: 0.8rem;
    min-height: 36px;
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 12px;
    font-weight: bold;
  }

  .articles__item {
    cursor: pointer;
  }

  .articles__404 {
    background-color: var(--color-level-1);
    height: 240px;
    border-radius: var(--border-radius);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 24px;
  }

  .overlay__article-manage {
    width: 100%;
    display: flex;
    flex-direction: column;
  }

  .overlay__item {
    height: 64px;
    width: 100%;
    font-size: 1rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .overlay__item:hover {
    background-color: var(--color-hover);
  }
</style>
