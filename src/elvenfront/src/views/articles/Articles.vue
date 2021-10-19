<template>
  <div class="articles__container">
    <div class="articles__create">
      <RouterLink :to="{ name: 'ArticleCreate' }">create</RouterLink>
    </div>
    <div class="articles__toolbar">
      <div class="articles__show">
        <div
          class="articles__item articles__show-published"
          v-if="show === 'published'"
          v-on:click="getArticles(undefined, 'drafts')"
        >show published</div>
        <div
          class="articles__item articles__show-drafts"
          v-if="show === 'drafts'"
          v-on:click="getArticles(undefined, 'published')"
        >show drafts</div>
      </div>
      <div class="articles__sort-by-date">
        <div
          class="articles__item"
          v-if="sortFirst === 'newest'"
          v-on:click="setSortDate('oldest')"
        >newest</div>
        <div
          class="articles__item"
          v-if="sortFirst === 'oldest'"
          v-on:click="setSortDate('newest')"
        >oldest</div>
      </div>
      <div class="articles__sort-by">
        <div
          class="articles__item"
          v-if="sortBy === 'updated'"
          v-on:click="setSort('published')"
        >by updated date</div>
        <div
          class="articles__item"
          v-if="sortBy === 'published'"
          v-on:click="setSort('created')"
        >by published date</div>
        <div
          class="articles__item"
          v-if="sortBy === 'created'"
          v-on:click="setSort('updated')"
        >by created date</div>
      </div>
    </div>

    <ArticlesList :articles="articles" @selected="selectArticle($emit)"></ArticlesList>

    <div class="articles__404" v-if="isArticlesLoaded && articles.length < 1">
      <div>no articles :(</div>
    </div>

    <Pagination
      :total-pages="totalPages"
      :current-page="currentPage"
      @changed="getArticles($event)"
    ></Pagination>
  </div>

  <Overlay v-bind:active="isToolsOverlayActive" v-on:deactivated="isToolsOverlayActive = false">
    <!-- tools -->
    <div class="overlay__article-manage" v-if="isToolsOverlayActive">
      <div
        class="overlay__item make__draft"
        v-if="selectedArticle && selectedArticle.is_published"
        v-on:click="makeDraftArticle(selectedArticle)"
      >make a draft</div>
      <div class="overlay__item publish" v-else v-on:click="publishArticle(selectedArticle)">publish</div>
      <div class="overlay__item edit" v-on:click="editArticle(selectedArticle)">edit</div>
      <div class="overlay__item delete" v-on:click="deleteArticle(selectedArticle)">delete</div>
    </div>
  </Overlay>
</template>

<script setup lang="ts">
import { onMounted, ref, Ref } from "@vue/runtime-core"
import { useRouter } from "vue-router"
import { IArticle } from "@/types/article"
import { IMeta } from "@/types/global"
import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter"
import Overlay from "@/components/ui/Overlay.vue"
import Pagination from "@/components/ui/Pagination.vue"
import ArticlesList from "@/components/parts/ArticlesList.vue"

const router = useRouter()
// service
const isArticlesLoaded = ref(false)
const isToolsOverlayActive = ref(false)
const selectedArticle: Ref<IArticle | null> = ref(null)
// articles
const articles: Ref<Array<IArticle>> = ref([])
const articlesMeta: Ref<IMeta | null> = ref(null)
// params
const show = ref('published')
const sortBy = ref('updated')
const sortFirst = ref('newest')
// pagination
const page = ref(1)
const perPage = ref(1)
const totalPages = ref(1)
const currentPage = ref(1)

onMounted(async () => {
  await getArticles()
})

async function getArticles(pageA = page.value, showA = show.value, sortByA = sortBy.value, sortFirstA = sortFirst.value) {
  page.value = pageA
  show.value = showA
  sortBy.value = sortByA
  sortFirst.value = sortFirstA
  isArticlesLoaded.value = false
  ArticleAdapter.getArticles(pageA, showA, sortByA, sortFirstA)
    .then(async result => {
      articles.value = result.data
      articlesMeta.value = result.meta
      perPage.value = articlesMeta.value.per_page
      currentPage.value = articlesMeta.value.current_page
      totalPages.value = articlesMeta.value.total_pages
      isArticlesLoaded.value = true
    })
}

async function editArticle(article) {
  await router.push({ name: 'ArticleCreate', params: { id: article.id } })
}

async function deleteArticle(article) {
  const isDelete = confirm('Удалить запись?')
  if (isDelete) {
    await ArticleAdapter.deleteArticle(article.id)
    deleteArticleFromArray(article)
    isToolsOverlayActive.value = false
    refreshArticles()
  }
}

async function publishArticle(article) {
  await ArticleAdapter.publishArticle(article)
    .then(() => {
      deleteArticleFromArray(article)
      isToolsOverlayActive.value = false
    })
  refreshArticles()
}

async function makeDraftArticle(article) {
  await ArticleAdapter.makeDraftArticle(article)
    .then(() => {
      deleteArticleFromArray(article)
      isToolsOverlayActive.value = false
    })
  refreshArticles()
}

async function refreshArticles() {
  // refresh is need when for ex. you deleted all articles on current page
  // and we need to check, is data on current page exists?
  // if page > 1 and no data, we moving back (currentPage--) and get new articles
  let notArticles = isArticlesLoaded.value && articles.value.length < 1
  console.log(notArticles)
  if (notArticles) { // no articles in current page
    while (notArticles) {
      // moving back until the pages ends or data appears
      currentPage.value--
      await getArticles()
      if (currentPage.value <= 1) {
        break
      }
      notArticles = isArticlesLoaded.value && articles.value.length < 1
    }
  }
}


async function setSort(sort) {
  sortBy.value = sort
  page.value = 1
  await getArticles()
}
async function setSortDate(age = 'newest') {
  sortFirst.value = age
  page.value = 1
  await getArticles()
}

function deleteArticleFromArray(article) {
  const index = articles.value.indexOf(article)
  articles.value.splice(index, 1)
  return true
}

function selectArticle(article) {
  isToolsOverlayActive.value = true
  selectedArticle.value = article
}

</script>

<style scoped>
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