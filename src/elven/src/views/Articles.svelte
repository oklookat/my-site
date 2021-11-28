<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import type { Article, By, Params, Show, Start } from "@/types/article";
  import type { Meta } from "@/types/global";
  import ArticleAdapter from "@/adapters/ArticleAdapter";
  import Overlay from "@/components/ui/Overlay.svelte";
  import Pagination from "@/components/ui/Pagination.svelte";
  import CArticle from "@/components/parts/Article.svelte";

  let loaded = false;
  let toolsOverlay = false;
  let selected: Article;
  let articles: Array<Article> = [];
  let meta: Meta;
  let params: Params = {
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

  async function getAll(p?: Params) {
    p = !p ? params : p;
    if (p.page < 1) {
      p.page = 1;
    }
    loaded = false;
    try {
      const result = await ArticleAdapter.getAll(params);
      articles = result.data;
      meta = result.meta;
      loaded = true;
    } catch (err) {}
  }

  async function edit(id: string) {
    await push(`/articles/create/${id}`);
  }

  function select(a: Article) {
    toolsOverlay = true;
    selected = a;
  }

  async function publish(id: string) {
    await ArticleAdapter.publish(id).then(() => {
      deleteFromArray(id);
      toolsOverlay = false;
    });
    await refresh();
  }

  async function toDrafts(id: string) {
    await ArticleAdapter.makeDraft(id).then(() => {
      deleteFromArray(id);
      toolsOverlay = false;
    });
    await refresh();
  }

  async function deleteArticle(id: string) {
    const isDelete = confirm("Delete article?");
    if (isDelete) {
      await ArticleAdapter.delete(id);
      await deleteFromArray(id);
      toolsOverlay = false;
    }
  }

  async function deleteFromArray(id: string) {
    articles = articles.filter((a) => a.id !== id);
    await refresh();
  }

  async function setBy(by: By = "published") {
    params.by = by;
    params.page = 1;
    await getAll();
  }

  async function setStart(start: Start = "newest") {
    params.start = start;
    params.page = 1;
    await getAll();
  }

  async function setShow(show: Show) {
    params.show = show;
    params.page = 1;
    await getAll();
  }

  async function onPageChanged(page: number) {
    params.page = page;
    await getAll();
  }

  async function refresh() {
    let noArticles =
      loaded && (articles.length < 1 || articles.length < meta.per_page);
    while (noArticles) {
      params.page--;
      await getAll();
      if (params.page <= 1) {
        break;
      }
    }
  }
</script>

<div class="articles">
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
      totalPages={meta.total_pages}
      currentPage={meta.current_page}
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
    &__toolbar {
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
    &__create {
      background-color: var(--color-level-1);
      height: 36px;
      width: 100%;
      display: flex;
      flex-direction: row;
      > a {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        :hover {
          background-color: var(--color-hover);
        }
      }
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
