<template>
  <div class="container">
    <Header title="записи"></Header>
    <div class="content">
      <div class="articles-list">
        <article class="article" v-for="article in articles" :key="article.id" v-on:click="selectArticle(article)">
          <div class="article-item article-title">{{ article.title }}</div>
          <div class="article-item article-updated-at">ред. {{ convertDateWrap(article.updated_at) }}</div>
        </article>
      </div>
      <div class="articles-404" v-if="isArticlesLoaded && articles.length < 1">
        <div>
          Нет записей. Но вы можете
          <RouterLink class="articles-404-link" :to="{name: 'ArticleCreate'}">создать новую</RouterLink>.
        </div>
      </div>
    </div>

    <UIOverlay v-bind:active="isOverlayActive" v-on:deactivated="isOverlayActive = false">
      <div class="article-tools">
        <div class="item article-action">
          <div class="article-make-draft" v-if="selectedArticle.is_published">Сделать черновиком</div>
          <div class="article-publish" v-else>Опубликовать</div>
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
import ElvenTools from "@/common/tools/ElvenTools";

export default defineComponent({
  name: 'Articles',
  components: {UIOverlay, Header},
  data(){
    return{
      isOverlayActive: false,
      articles: [],
      articlesMeta: [],
      selectedArticle: undefined,
      elvenTools: undefined,
      isArticlesLoaded: false,
    }
  },
  async mounted() {
    await this.getArticles()
  },
  methods: {
    async getArticles(page = '1'){
      this.isArticlesLoaded = false;
      await ArticleAdapter.getArticles(page)
        .then(result =>{
          this.articles = result.data
          this.articlesMeta = result.meta
          this.isArticlesLoaded = true;
        })
    },
    selectArticle(article){
      this.isOverlayActive = true
      this.selectedArticle = article
    },
    async deleteArticle(article){
      const isDelete = confirm('Удалить запись?')
      if(isDelete){
        await ArticleAdapter.deleteArticle(article.id)
        const index = this.articles.indexOf(article)
        this.articles.splice(index, 1)
        this.isOverlayActive = false
      }
    },
    async editArticle(article){
      await this.$router.push({name: 'ArticleCreate', params: {id: article.id}})
    },
    convertDateWrap(date){
      return ElvenTools.convertDate(date)
    }
  }
})
</script>

<style scoped lang="scss">
.container{

}
.articles-list{
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  min-height: 42px;
  gap: 12px;

}
.article{
  border-radius: 6px;
  background-color: var(--color-level-2);
  width: 100%;
  height: 100%;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  padding-bottom: 12px;
}
.article:hover{
  background-color: var(--color-hover);
}
.article-item{
  font-size: 1.2rem;
  line-height: 1.5rem;
  margin-top: 12px;
  margin-left: 12px;
  margin-right: 12px;
}
.article-updated-at{
  font-size: 0.9rem;
}
.article-tools{
  width: 100%;
  display: flex;
  flex-direction: column;
}
.item{
  display: flex;
  align-items: center;
  justify-content: center;
  height: 64px;
  width: 100%;
}
.articles-404{
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
}
.articles-404 a{
  text-decoration: underline;
}



@media screen and (min-width: 1023px) {
  .articles-list{

  }
}
</style>