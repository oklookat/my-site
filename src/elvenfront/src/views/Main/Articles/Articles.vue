<template>
  <div class="container">
    <Header></Header>
    <div class="content">
      <div class="articles-types">
        <div class="articles-published" v-on:click="getArticles(undefined, 'published')"
             :class="{'active': show === 'published'}">
          Статьи
        </div>
        <div class="articles-drafts" v-on:click="getArticles(undefined, 'drafts')"
             :class="{'active': show === 'drafts'}">
          Черновики
        </div>
      </div>
      <div class="articles-sort-by">
        <div class="articles-sort-by-link">
          <div v-if="sortBy === 'created'" v-on:click="isSortOverlayActive = !isSortOverlayActive">дата создания</div>
          <div v-if="sortBy === 'updated'" v-on:click="isSortOverlayActive = !isSortOverlayActive">последнее
            редактирование
          </div>
          <div v-if="sortBy === 'published'" v-on:click="isSortOverlayActive = !isSortOverlayActive">дата публикации
          </div>
          <div class="the-underline"></div>
        </div>
        <div class="articles-sort-by-link">
          <div v-if="sortFirst === 'newest'"
               v-on:click="setSortDate('oldest')">
            новые
          </div>
          <div v-if="sortFirst === 'oldest'"
               v-on:click="setSortDate('newest')">
            старые
          </div>
          <div class="the-underline"></div>
        </div>
      </div>

      <div class="articles-list">
        <article class="article" v-for="article in articles" :key="article.id" v-on:click="selectArticle(article)">
          <div class="article-meta">
            <div class="article-item article-updated-at" v-if="!article.is_published">
              {{ convertDateWrap(article.updated_at) }}
            </div>
            <div class="article-item article-published-at" v-else>
              {{ convertDateWrap(article.published_at) }}
            </div>
            <div class="article-item article-is-draft" v-if="!article.is_published">Черновик</div>
          </div>
          <div class="article-main">
            <div class="article-item article-title">{{ article.title }}</div>
            <div class="article-item article-content-preview">
              {{ article.content }}
            </div>
          </div>
        </article>
      </div>
      <div class="articles-404" v-if="isArticlesLoaded && articles.length < 1">
        <div class="articles-404-1">Нет записей.</div>
        <div class="articles-404-2">
          <RouterLink class="articles-404-link" :to="{name: 'ArticleCreate'}">
            Создать новую?
            <div class="the-underline"></div>
          </RouterLink>
        </div>
      </div>
    </div>

    <UIOverlay v-bind:active="isToolsOverlayActive" v-on:deactivated="isToolsOverlayActive = false">
      <div class="overlay-article-tools">
        <div class="ov-item article-action">
          <div class="article-make-draft" v-if="selectedArticle.is_published"
               v-on:click="makeDraftArticle(selectedArticle)">
            Сделать черновиком
          </div>
          <div class="article-publish" v-else v-on:click="publishArticle(selectedArticle)">Опубликовать</div>
        </div>
        <div class="ov-item article-edit" v-on:click="editArticle(selectedArticle)">Редактировать</div>
        <div class="ov-item article-delete" v-on:click="deleteArticle(selectedArticle)">Удалить</div>
      </div>
    </UIOverlay>

    <UIOverlay v-bind:active="isSortOverlayActive" v-on:deactivated="isSortOverlayActive = false">
      <div class="overlay-article-sort">
        <div class="ov-item asb-sort-by-created"
             :class="{'active': sortBy === 'created'}"
             v-on:click="setSort('created')">Дата создания
        </div>
        <div class="ov-item asb-sort-by-updated"
             :class="{'active': sortBy === 'updated'}"
             v-on:click="setSort('updated')">Дата последнего редактирования
        </div>
        <div class="ov-item asb-sort-by-published"
             :class="{'active': sortBy === 'published'}"
             v-on:click="setSort('published')">Дата публикации
        </div>
      </div>
    </UIOverlay>
  </div>
</template>

<script>
import {defineComponent} from "vue";
import Header from "@/components/Header/Header";
import ArticleAdapter from "@/common/adapters/Main/ArticleAdapter";
import UIOverlay from "@/components/_UI/UIOverlay";
import ElvenDates from "@/common/tools/ElvenDates";

export default defineComponent({
  name: 'Articles',
  components: {UIOverlay, Header},
  data() {
    return {
      isArticlesLoaded: false,
      show: 'published', // drafts or published
      isToolsOverlayActive: false, // actions overlay
      articles: [],
      articlesMeta: [],
      selectedArticle: undefined,
      currentPage: 1,

      sortBy: 'updated', // see backend docs for more
      sortFirst: 'newest',
      isSortOverlayActive: false,
    }
  },
  async mounted() {
    await this.getArticles()
  },
  methods: {
    async getArticles(page = this.currentPage, show = this.show, sortBy = this.sortBy, sortFirst = this.sortFirst) {
      this.show = show
      this.isArticlesLoaded = false;
      await ArticleAdapter.getArticles(page, show, sortBy, sortFirst)
          .then(result => {
            this.articles = result.data
            this.articlesMeta = result.meta
            this.isArticlesLoaded = true;
          })
    },
    async deleteArticle(article) {
      const isDelete = confirm('Удалить запись?')
      if (isDelete) {
        await ArticleAdapter.deleteArticle(article.id)
        const index = this.articles.indexOf(article)
        this.articles.splice(index, 1)
        this.isToolsOverlayActive = false
      }
    },
    async editArticle(article) {
      await this.$router.push({name: 'ArticleCreate', params: {id: article.id}})
    },
    async publishArticle(article) {
      await ArticleAdapter.publishArticle(article)
    },
    async makeDraftArticle(article) {
      await ArticleAdapter.makeDraftArticle(article)
    },
    selectArticle(article) {
      this.isToolsOverlayActive = true
      this.selectedArticle = article
    },
    convertDateWrap(date) {
      return ElvenDates.convert(date)
    },
    async setSort(sort) {
      this.sortBy = sort
      await this.getArticles()
      this.isSortOverlayActive = false
    },
    async setSortDate(age = 'newest'){
      this.sortFirst = age
      await this.getArticles()
    }
  }
})
</script>

<style scoped>
.container {

}

.content {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.articles-types {
  background-color: var(--color-level-1);
  border-radius: 6px;
  height: 42px;
  width: 100%;
  display: flex;
  flex-direction: row;
}

.articles-types div {
  border-radius: inherit;
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

.articles-sort-by {
  background-color: var(--color-level-1);
  color: var(--color-text-inactive);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding-left: 12px;
  font-size: 0.8rem;
  min-height: 36px;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
}

.articles-sort-by-link
{
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
  border-radius: 6px;
  background-color: var(--color-level-1);
  cursor: pointer;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding-bottom: 12px;
  gap: 8px;
}

.article-meta {
  display: flex;
  flex-direction: row;
  color: var(--color-text-inactive);
}

.article-item {
  font-size: 1.1rem;
  line-height: 1.5rem;
  margin-top: 8px;
  margin-left: 12px;
  margin-right: 12px;
}

.article-title {
  font-size: 1.6rem;
}

.article-content-preview {
  height: 100%;
}

.article-updated-at,
.article-published-at,
.article-is-draft {
  font-size: 0.9rem;
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

.articles-404 {
  display: grid;
  grid-template-rows: repeat(2, 1fr);
  grid-template-columns: 1fr;
  justify-items: center;
  grid-gap: 24px;
}


@media screen and (min-width: 512px) {
  .article-content-preview {
    width: 400px;
  }
}

@media screen and (min-width: 1024px) {
  .article-content-preview {
    width: 412px;
  }
}
</style>