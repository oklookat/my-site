<script lang="ts">
  import type { IArticle } from "@/types/article";
  import type { IMeta } from "@/types/global";
  import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter";
  import Overlay from "@/components/ui/Overlay.svelte";
  import Pagination from "@/components/ui/Pagination.svelte";
  import ArticlesList from "@/components/parts/ArticlesList.svelte";
  import { onMount } from "svelte";

  // service
  let isArticlesLoaded = false;
  let isToolsOverlayActive = false;
  let selectedArticle: IArticle | null = null;
  // articles
  let articles: Array<IArticle> = [];
  let articlesMeta: IMeta | null = null;
  // params
  let show = "published";
  let sortBy = "updated";
  let sortFirst = "newest";
  // pagination
  let page = 1;
  let perPage = 1;
  let totalPages = 1;
  let currentPage = 1;

  onMount(async () => {
    await getArticles();
  });

  async function getArticles(
    pageA = page,
    showA = show,
    sortByA = sortBy,
    sortFirstA = sortFirst
  ) {
    page = pageA;
    show = showA;
    sortBy = sortByA;
    sortFirst = sortFirstA;
    isArticlesLoaded = false;
    ArticleAdapter.getArticles(pageA, showA, sortByA, sortFirstA).then(
      async (result) => {
        articles = result.data;
        articlesMeta = result.meta;
        perPage = articlesMeta.per_page;
        currentPage = articlesMeta.current_page;
        totalPages = articlesMeta.total_pages;
        isArticlesLoaded = true;
      }
    );
  }

  async function editArticle(article) {
    //await router.push({ name: 'ArticleCreate', params: { id: article.id } })
  }

  async function deleteArticle(article) {
    const isDelete = confirm("delete article?");
    if (isDelete) {
      await ArticleAdapter.deleteArticle(article.id);
      deleteArticleFromArray(article);
      isToolsOverlayActive = false;
      refreshArticles();
    }
  }

  async function publishArticle(article) {
    await ArticleAdapter.publishArticle(article).then(() => {
      deleteArticleFromArray(article);
      isToolsOverlayActive = false;
    });
    refreshArticles();
  }

  async function makeDraftArticle(article) {
    await ArticleAdapter.makeDraftArticle(article).then(() => {
      deleteArticleFromArray(article);
      isToolsOverlayActive = false;
    });
    refreshArticles();
  }

  async function refreshArticles() {
    // refresh is need when for ex. you deleted all articles on current page
    // and we need to check, is data on current page exists?
    // if page > 1 and no data, we moving back (currentPage--) and get new articles
    let notArticles = isArticlesLoaded && articles.length < 1;
    console.log(notArticles);
    if (notArticles) {
      // no articles in current page
      while (notArticles) {
        // moving back until the pages ends or data appears
        currentPage--;
        await getArticles();
        if (currentPage <= 1) {
          break;
        }
        notArticles = isArticlesLoaded && articles.length < 1;
      }
    }
  }

  async function setSort(sort) {
    sortBy = sort;
    page = 1;
    await getArticles();
  }
  async function setSortDate(age = "newest") {
    sortFirst = age;
    page = 1;
    await getArticles();
  }

  function deleteArticleFromArray(article) {
    const index = articles.indexOf(article);
    articles.splice(index, 1);
    return true;
  }

  function selectArticle(article) {
    isToolsOverlayActive = true;
    selectedArticle = article;
  }
</script>

<main>
  <div class="articles__container">
    <div class="articles__create">
      <!-- <RouterLink :to="{ name: 'ArticleCreate' }">create</RouterLink> -->
    </div>
    <div class="articles__toolbar">
      <div class="articles__show">
        {#if show === "published"}
          <div
            class="articles__item articles__show-published"
            on:click={() => getArticles(undefined, "drafts")}
          >
            show published
          </div>
        {/if}
        {#if show === "drafts"}
          <div
            class="articles__item articles__show-drafts"
            on:click={() => getArticles(undefined, "published")}
          >
            show drafts
          </div>
        {/if}
      </div>
      <div class="articles__sort-by-date">
        {#if sortFirst === "newest"}
          <div class="articles__item" on:click={() => setSortDate("oldest")}>
            newest
          </div>
        {/if}
        {#if sortFirst === "oldest"}
          <div class="articles__item" on:click={() => setSortDate("newest")}>
            oldest
          </div>
        {/if}
      </div>
      <div class="articles__sort-by">
        {#if sortBy === "updated"}
          <div class="articles__item" on:click={() => setSort("published")}>
            by updated date
          </div>
        {/if}
        {#if sortBy === "published"}
          <div
            class="articles__item"
            v-if="sortBy === 'published'"
            on:click={() => setSort("created")}
          >
            by published date
          </div>
        {/if}
        {#if sortBy === "created"}
          <div class="articles__item" on:click={() => setSort("updated")}>
            by created date
          </div>
        {/if}
      </div>

      <ArticlesList articles="articles" on:selected={selectArticle} />

      <div class="articles__404" v-if="isArticlesLoaded && articles.length < 1">
        <div>no articles :(</div>
      </div>

      <Pagination
        total-pages="totalPages"
        current-page="currentPage"
        on:changed={getArticles}
      />
    </div>

    <Overlay
      bind:active={isToolsOverlayActive}
      on:deactivated={() => (isToolsOverlayActive = false)}
    >
      <!-- tools -->
      <div class="overlay__article-manage" v-if="isToolsOverlayActive">
        {#if selectedArticle && selectedArticle.is_published}
          <div
            class="overlay__item make__draft"
            on:click={() => makeDraftArticle(selectedArticle)}
          >
            make a draft
          </div>
        {:else}
          <div
            class="overlay__item publish"
            v-on:click="publishArticle(selectedArticle)"
          >
            publish
          </div>
        {/if}
        <div
          class="overlay__item edit"
          on:click={() => editArticle(selectedArticle)}
        >
          edit
        </div>
        <div
          class="overlay__item delete"
          on:click={() => deleteArticle(selectedArticle)}
        >
          delete
        </div>
      </div>
    </Overlay>
  </div>
</main>

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

  .overlay__article-manage,
  .overlay__articles-sort {
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
