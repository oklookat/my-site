<template>
  <div class="container">
    <Header></Header>
    <div class="content">
      <div class="articles-types">
        <div class="articles-published" v-on:click="getPublished()"
             :class="{'articles-types-active': !isShowDrafts}">
          Статьи
        </div>
        <div class="articles-drafts" v-on:click="getDrafts()"
             :class="{'articles-types-active': isShowDrafts}">
          Черновики
        </div>
      </div>
      <!--      <div class="articles-sort-by">дата </div>-->

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
          <div class="article-item article-title">{{ article.title }}</div>
          <div class="article-item article-content-preview">
            {{ showPreviewContent(article.content) }}
          </div>
        </article>
      </div>
      <div class="articles-404" v-if="isArticlesLoaded && articles.length < 1">
        <div class="articles-404-1">Нет записей.</div>
        <div class="articles-404-2">
          Но вы можете
          <RouterLink class="articles-404-link" :to="{name: 'ArticleCreate'}">создать новую</RouterLink>.
        </div>
      </div>
    </div>

    <UIOverlay v-bind:active="isOverlayActive" v-on:deactivated="isOverlayActive = false">
      <div class="article-tools">
        <div class="item article-action">
          <div class="article-make-draft" v-if="selectedArticle.is_published"
               v-on:click="makeDraftArticle(selectedArticle)">Сделать черновиком
          </div>
          <div class="article-publish" v-else v-on:click="publishArticle(selectedArticle)">Опубликовать</div>
        </div>
        <div class="item article-edit" v-on:click="editArticle(selectedArticle)">Редактировать</div>
        <div class="item article-delete" v-on:click="deleteArticle(selectedArticle)">Удалить</div>
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
      isShowDrafts: false,

      isOverlayActive: false,
      articles: [],
      articlesMeta: [],
      selectedArticle: undefined,
      elvenTools: undefined,
      isArticlesLoaded: false,
    }
  },
  async mounted() {
    await this.getPublished()
  },
  methods: {
    async getPublished(page = '1'){
      this.isArticlesLoaded = false;
      this.isShowDrafts = false;
      await ArticleAdapter.getPublished(page)
          .then(result => {
            this.articles = result.data
            this.articlesMeta = result.meta
            this.isArticlesLoaded = true;
          })
    },
    async getDrafts(page = '1'){
      this.isArticlesLoaded = false;
      this.isShowDrafts = true;
      await ArticleAdapter.getDrafts(page)
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
        this.isOverlayActive = false
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
      this.isOverlayActive = true
      this.selectedArticle = article
    },
    convertDateWrap(date) {
      return ElvenDates.convert(date)
    },
    showPreviewContent(content){
      content = JSON.parse(content)
      let text = content.blocks[0].data.text
      if(text.length > 408){
        text = text.slice(0, 408) + '...'
      }
      return text
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
  height: 42px;
  width: 100%;
  display: flex;
  flex-direction: row;
  gap: 24px;
}
.articles-types div {
  border-radius: 6px;
  width: 100%;
  height: 100%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
.articles-types div:hover{
  background-color: var(--color-hover);
}
.articles-types-active{
  background-color: var(--color-hover);
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
  background-color: var(--color-level-2);
  width: 100%;
  height: 100%;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  padding-bottom: 12px;
  gap: 8px;
}

.article:hover {
  background-color: var(--color-hover);
}
.article-meta{
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
.article-title{
  font-size: 1.6rem;
}
.article-content-preview{
  height: 100%;
}
.article-updated-at,
.article-published-at,
.article-is-draft {
  font-size: 0.9rem;
}


.article-tools {
  width: 100%;
  display: flex;
  flex-direction: column;
}

.item {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 64px;
  width: 100%;
}

.articles-404 {
  display: grid;
  grid-template-rows: repeat(2, 1fr);
  grid-template-columns: 1fr;
  justify-items: center;
  grid-gap: 24px;
}
.articles-404 a {
  text-decoration: underline;
}

@media screen and (min-width: 512px) {
  .article-content-preview{
    width: 400px;
  }
}
@media screen and (min-width: 1024px) {
  .article-content-preview{
    width: 412px;
  }
}
</style>