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

      <div class="elven-pagination" v-if="isPagination">
        <div class="paginator">
          <div class="prev-page">
            <div class="prev-page-butt"
                 v-if="currentPage !== 1" v-on:click="getArticles(currentPage - 1)">
              назад
            </div>
          </div>
          <div class="pages-numbers">
            <div class="page-number"
                 :class="{'active': currentPage === page}"
                 v-bind:key="page"
                 v-for="page in pagesNumbers"
                 v-on:click="getArticles(page)">
              {{ page }}
            </div>
          </div>
          <div class="next-page">
            <div class="next-page-butt"
                 v-if="currentPage < pagesCount" v-on:click="getArticles(currentPage + 1)">
              вперед
            </div>
          </div>
        </div>
      </div>

    </div>


    <UIOverlay v-bind:active="isToolsOverlayActive" v-on:deactivated="isToolsOverlayActive = false">
      <div class="overlay-article-tools">
        <div class="ov-item article-make-draft" v-if="selectedArticle.is_published"
             v-on:click="makeDraftArticle(selectedArticle)">
          Сделать черновиком
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

      // PAGINATION START //
      isPagination: false,
      pagesCount: 1,
      currentPage: 1,
      perPage: 1,
      pagesNumbers: [],
      // PAGINATION END //

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
          .then(async result => {
            this.articles = result.data
            this.articlesMeta = result.meta
            this.pagesCount = this.articlesMeta.total
            this.perPage = this.articlesMeta.per_page
            this.currentPage = this.articlesMeta.current_page
            this.isPagination = this.pagesCount > 1
            this.isArticlesLoaded = true
          })
      if(this.isPagination){
        this.generatePageNumbers()
      }
    },
    async deleteArticle(article) {
      const isDelete = confirm('Удалить запись?')
      if (isDelete) {
        await ArticleAdapter.deleteArticle(article.id)
        this.deleteArticleFromArray(article)
        this.isToolsOverlayActive = false
      }
    },
    async editArticle(article) {
      await this.$router.push({name: 'ArticleCreate', params: {id: article.id}})
    },
    async publishArticle(article) {
      await ArticleAdapter.publishArticle(article)
          .then(() => {
            this.deleteArticleFromArray(article)
            this.isToolsOverlayActive = false
          })
    },
    async makeDraftArticle(article) {
      await ArticleAdapter.makeDraftArticle(article)
          .then(() => {
            this.deleteArticleFromArray(article)
            this.isToolsOverlayActive = false
          })
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
    generatePageNumbers(){
      // https://stackoverflow.com/a/35109214
      const currentPage = this.currentPage
      const totalPages = this.pagesCount
      const pageSize = 10 // maximum elements in paginator (UI, not backend)
      let startPage, endPage
      if (totalPages <= pageSize) {
        startPage = 1
        endPage = totalPages
      } else {
        if (currentPage <= 6) {
          startPage = 1
          endPage = 10
        } else if (currentPage + 4 >= totalPages) {
          startPage = totalPages - 9
          endPage = totalPages
        } else {
          startPage = currentPage - 5
          endPage = currentPage + 4
        }
      }
      let pages = []
      for(let i = 0; startPage < endPage + 1; i++){
        pages[i] = startPage
        startPage++
      }
      this.pagesNumbers = pages
    }
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
  border-radius: 6px;
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
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 24px;
}

.elven-pagination {
  border-radius: 8px;
  background-color: var(--color-level-1);
  height: 64px;
}

.paginator {
  border-radius: inherit;
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: row;
}

.next-page,
.prev-page {
  border-radius: inherit;
  width: 25%;
}

.next-page-butt,
.prev-page-butt {
  border-radius: inherit;
  cursor: pointer;
  width: 100%;
  height: 100%;
}

.next-page-butt:hover,
.prev-page-butt:hover {
  background-color: var(--color-hover);
}

.next-page-butt,
.prev-page-butt,
.pages-numbers {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.pages-numbers {
  width: 50%;
}

.pages-numbers > .page-number {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}

.page-number:hover {
  background-color: var(--color-hover);
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