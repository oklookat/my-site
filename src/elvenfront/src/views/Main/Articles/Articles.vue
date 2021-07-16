<template>
  <div class="container">
    <Header title="записи"></Header>
    <div class="content">
      <div class="articles-list">
        <article class="article" v-for="article in articles" :key="article.id" v-on:click="selectArticle(article)">
          <div class="article-item article-title">{{ article.title }}</div>
          <div class="article-item article-updated-at">Изменена: {{ convertDateWrap(article.updated_at) }}</div>
        </article>
      </div>
      <h1>ПРИ СКЛОНЕНИИ ДАТ КАКАЯ-ТО БЕДА. ЕСЛИ ОБНОВЛЯТЬ СТРАНИЦУ, ДАТЫ СТРАННО СКЛОНЯЮТСЯ.</h1>
      <div class="articles-404" v-if="articles.length < 1">
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
    }
  },
  async mounted() {
    await this.getArticles()
  },
  methods: {
    async getArticles(page = '1'){
      await ArticleAdapter.getArticles(page)
        .then(result =>{
          console.log(result)
          this.articles = result.data
          this.articlesMeta = result.meta
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
.content{
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: 42px 1fr;
  grid-gap: 12px;
}
.articles-list{
  height: 100%;
  width: 100%;
  display: grid;
  grid-template-columns: 1fr;
  grid-auto-rows: minmax(64px, 1fr);
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
}
.article:hover{
  background-color: var(--color-hover);
}
.article-item{
  font-size: 1.4rem;
  margin: 4px;
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