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

      <div class="articles-list" v-if="articles.length > 0">
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
        <div class="articles-404-1">Нет записей :(</div>
        <div class="articles-404-2">
          <RouterLink class="articles-404-link" :to="{name: 'ArticleCreate'}">
            Создать новую?
            <div class="the-underline"></div>
          </RouterLink>
        </div>
      </div>

      <UIPagination
          :total-pages="totalPages"
          :current-page="currentPage"
          v-on:page-changed="getArticles($event)">
      </UIPagination>
    </div>


    <UIOverlay v-bind:active="isToolsOverlayActive" v-on:deactivated="isToolsOverlayActive = false">
      <div class="overlay-article-tools">
        <div class="ov-item article-make-draft" v-if="selectedArticle.is_published"
             v-on:click="makeDraftArticle(selectedArticle)">
          В черновики
        </div>
        <div class="ov-item article-publish" v-else v-on:click="publishArticle(selectedArticle)">Опубликовать</div>
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
import UIPagination from "@/components/_UI/UIPagination";

export default defineComponent({
  name: 'Articles',
  components: {UIPagination, UIOverlay, Header},
  data() {
    return {
      isArticlesLoaded: false,
      isToolsOverlayActive: false,
      isSortOverlayActive: false,
      show: 'published', // see backend docs for more
      sortBy: 'updated', // see backend docs for more
      sortFirst: 'newest',
      articles: [],
      articlesMeta: [],
      selectedArticle: undefined,
      totalPages: 1,
      currentPage: 1,
    }
  },
  async mounted() {
    await this.getArticles()
  },
  watch: {},
  methods: {
    async getArticles(page = this.currentPage, show = this.show, sortBy = this.sortBy, sortFirst = this.sortFirst) {
      if (show !== this.show || this.currentPage < 1) {
        this.currentPage = 1
        page = this.currentPage
      }
      this.show = show
      this.isArticlesLoaded = false;
      await ArticleAdapter.getArticles(page, show, sortBy, sortFirst)
          .then(async result => {
            this.articles = result.data
            this.articlesMeta = result.meta
            this.perPage = this.articlesMeta.per_page
            this.currentPage = this.articlesMeta.current_page
            this.totalPages = Math.ceil(this.articlesMeta.total / this.articlesMeta.per_page)
            this.isArticlesLoaded = true
          })
    },
    async refreshArticles() {
      // refresh is need when for ex. you deleted all articles on current page
      // and we need to check, is data on current page exists?
      // if page > 1 and no data, we moving back (currentPage--) and get new articles
      let isTrueArticles = this.isArticlesLoaded && this.articles.length < 1
      if (isTrueArticles) { // no articles in current page
        while (isTrueArticles) {
          // moving back until the pages ends or data appears
          this.currentPage--
          await this.getArticles()
          if (this.currentPage <= 1) {
            break
          }
          isTrueArticles = this.isArticlesLoaded && this.articles.length < 1
        }
      }
    },
    async editArticle(article) {
      await this.$router.push({name: 'ArticleCreate', params: {id: article.id}})
    },
    async deleteArticle(article) {
      const isDelete = confirm('Удалить запись?')
      if (isDelete) {
        await ArticleAdapter.deleteArticle(article.id)
        this.deleteArticleFromArray(article)
        this.isToolsOverlayActive = false
        await this.refreshArticles()
      }
    },
    async publishArticle(article) {
      await ArticleAdapter.publishArticle(article)
          .then(() => {
            this.deleteArticleFromArray(article)
            this.isToolsOverlayActive = false
          })
      await this.refreshArticles()
    },
    async makeDraftArticle(article) {
      await ArticleAdapter.makeDraftArticle(article)
          .then(() => {
            this.deleteArticleFromArray(article)
            this.isToolsOverlayActive = false
          })
      await this.refreshArticles()
    },
    async setSort(sort) {
      this.sortBy = sort
      this.currentPage = 1
      await this.getArticles()
      this.isSortOverlayActive = false
    },
    async setSortDate(age = 'newest') {
      this.sortFirst = age
      this.currentPage = 1
      await this.getArticles()
    },

    // SERVICE START //
    deleteArticleFromArray(article) {
      const index = this.articles.indexOf(article)
      this.articles.splice(index, 1)
      return true
    },
    selectArticle(article) {
      this.isToolsOverlayActive = true
      this.selectedArticle = article
    },
    convertDateWrap(date) {
      return ElvenDates.convert(date)
    },
    // SERVICE END //
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
  border-radius: var(--border-radius);
  padding-left: 12px;
  font-size: 0.8rem;
  min-height: 36px;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
}

.articles-sort-by-link {
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