<template>
  <div class="articles-main">
    <div class="articles-create">
      <RouterLink :to="{ name: 'ArticleCreate' }">Создать</RouterLink>
    </div>
    <div class="articles-types">
      <div
        class="articles-published"
        v-on:click="getArticles(undefined, 'published')"
        :class="{ 'active': show === 'published' }"
      >Статьи</div>
      <div
        class="articles-drafts"
        v-on:click="getArticles(undefined, 'drafts')"
        :class="{ 'active': show === 'drafts' }"
      >Черновики</div>
    </div>
    <div class="articles-toolbar">
      <div class="articles-sort-by-link">
        <div
          v-if="sortBy === 'created'"
          v-on:click="isSortOverlayActive = !isSortOverlayActive"
        >создано</div>
        <div
          v-if="sortBy === 'updated'"
          v-on:click="isSortOverlayActive = !isSortOverlayActive"
        >изменено</div>
        <div
          v-if="sortBy === 'published'"
          v-on:click="isSortOverlayActive = !isSortOverlayActive"
        >опубликовано</div>
      </div>
      <div class="articles-sort-by-link">
        <div v-if="sortFirst === 'newest'" v-on:click="setSortDate('oldest')">новые</div>
        <div v-if="sortFirst === 'oldest'" v-on:click="setSortDate('newest')">старые</div>
      </div>
    </div>

    <div class="articles-list" v-if="articles.length > 0">
      <article
        class="article"
        v-for="article in articles"
        :key="article.id"
        v-on:click="selectArticle(article)"
      >
        <div class="article-meta">
          <div
            class="article-item article-updated-at"
            v-if="!article.is_published"
          >{{ convertDateWrap(article.updated_at) }}</div>
          <div
            class="article-item article-published-at"
            v-else
          >{{ convertDateWrap(article.published_at) }}</div>
          <div class="article-item article-is-draft" v-if="!article.is_published">Черновик</div>
        </div>
        <div class="article-main">
          <div class="article-item article-title">{{ article.title }}</div>
          <div
            v-if="article.content && article.content.blocks && article.content.blocks[0]"
            class="article-item article-content-preview"
          >
            <div v-html="article.content.blocks[0].data.text"></div>
          </div>
        </div>
      </article>
    </div>

    <div class="articles-404" v-if="isArticlesLoaded && articles.length < 1">
      <div class="articles-404-1">Нет записей :(</div>
      <div class="articles-404-2">
        <RouterLink class="articles-404-link" :to="{ name: 'ArticleCreate' }">Создать новую?</RouterLink>
      </div>
    </div>

    <UIPagination
      :total-pages="totalPages"
      :current-page="currentPage"
      v-on:changed="getArticles($event)"
    ></UIPagination>
  </div>

  <UIOverlay v-bind:active="isToolsOverlayActive" v-on:deactivated="isToolsOverlayActive = false">
    <div class="overlay-article-tools">
      <div
        class="ov-item article-make-draft"
        v-if="selectedArticle.is_published"
        v-on:click="makeDraftArticle(selectedArticle)"
      >В черновики</div>
      <div
        class="ov-item article-publish"
        v-else
        v-on:click="publishArticle(selectedArticle)"
      >Опубликовать</div>
      <div class="ov-item article-edit" v-on:click="editArticle(selectedArticle)">Редактировать</div>
      <div class="ov-item article-delete" v-on:click="deleteArticle(selectedArticle)">Удалить</div>
    </div>
  </UIOverlay>
  <UIOverlay v-bind:active="isSortOverlayActive" v-on:deactivated="isSortOverlayActive = false">
    <div class="overlay-article-sort">
      <div
        class="ov-item asb-sort-by-created"
        :class="{ 'active': sortBy === 'created' }"
        v-on:click="setSort('created')"
      >Дата создания</div>
      <div
        class="ov-item asb-sort-by-updated"
        :class="{ 'active': sortBy === 'updated' }"
        v-on:click="setSort('updated')"
      >Дата изменения</div>
      <div
        class="ov-item asb-sort-by-published"
        :class="{ 'active': sortBy === 'published' }"
        v-on:click="setSort('published')"
      >Дата публикации</div>
    </div>
  </UIOverlay>
</template>

<script setup lang="ts">
import { onMounted, ref, Ref } from "@vue/runtime-core"
import { useRouter, useRoute } from 'vue-router'
import { IArticle } from "@/types/article";
import { IMeta } from "@/types/global";
import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter";
import Dates from "@/common/tools/Dates";
import UIOverlay from "@/components/_UI/UIOverlay.vue";
import UIPagination from "@/components/_UI/UIPagination.vue";

const router = useRouter()
// service
const isArticlesLoaded = ref(false)
const isToolsOverlayActive = ref(false)
const isSortOverlayActive = ref(false)
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
  isSortOverlayActive.value = false
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

function convertDateWrap(date) {
  return Dates.convert(date)
}
</script>

<style scoped>
.articles-main {
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0px 0px 35px 9px rgba(34, 60, 80, 0.1);
}

.articles-types,
.articles-create {
  background-color: var(--color-level-1);
  height: 42px;
  width: 100%;
  display: flex;
  flex-direction: row;
}

.articles-create > a {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
}

.articles-create > a:hover {
  background-color: var(--color-hover);
}

.articles-types div {
  width: 100%;
  height: 100%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.articles-types div:hover {
  background-color: var(--color-hover);
}

.articles-toolbar {
  background-color: var(--color-level-1);
  color: var(--color-text-inactive);
  border: 1px solid var(--color-border);
  text-decoration: underline;
  padding-left: 12px;
  font-size: 0.8rem;
  min-height: 36px;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
  font-weight: bold;
  cursor: pointer;
}

.articles-list {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  min-height: 42px;
  gap: 12px;
}

.article {
  min-height: 164px;
  border-radius: var(--border-radius);
  background-color: var(--color-level-1);
  cursor: pointer;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding-bottom: 12px;
  gap: 8px;
}

.article-item {
  font-size: 1.1rem;
  line-height: 1.5rem;
  margin-top: 8px;
  margin-left: 12px;
  margin-right: 12px;
}

.article-meta {
  display: flex;
  flex-direction: row;
  color: var(--color-text-inactive);
}

.article-main {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.article-title {
  font-size: 1.5rem;
  line-height: 2rem;
  letter-spacing: 0.0099rem;
}

.article-content-preview {
  height: 100%;
}

.article-updated-at,
.article-published-at,
.article-is-draft {
  font-size: 0.9rem;
}

.articles-404 {
  background-color: var(--color-level-1);
  height: 240px;
  border-radius: var(--border-radius);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 24px;
}

.articles-404-link {
  text-decoration: underline;
}

.overlay-article-tools,
.overlay-article-sort {
  width: 100%;
  display: flex;
  flex-direction: column;
}

.ov-item {
  height: 64px;
  width: 100%;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ov-item:hover {
  background-color: var(--color-hover);
}

@media screen and (min-width: 512px) {
  .article-content-preview {
    width: 400px;
  }

  .article-main {
    width: 400px;
  }
}

@media screen and (min-width: 1024px) {
  .article-content-preview {
    width: 412px;
  }

  .article-main {
    width: 75%;
  }
}
</style>