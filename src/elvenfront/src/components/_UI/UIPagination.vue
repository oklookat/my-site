<template>
  <div class="ui-pagination" v-if="isPagination">
    <div class="ui-pagination-paginator">
      <div class="ui-pagination-prev-page">
        <div class="ui-pagination-prev-page-butt"
             v-if="currentPage_data !== 1" v-on:click="onPrevButton">
          назад
        </div>
      </div>
      <div class="ui-pagination-pages-input">
        <input class="ui-pagination-pages-input-num" type="text"
               placeholder="Номер страницы"
               v-model="inputPage_data"
               @input="onPageInput"/>
      </div>
      <div class="ui-pagination-next-page">
        <div class="ui-pagination-next-page-butt"
             v-if="currentPage_data < totalPages_data" v-on:click="onNextButton">
          вперед
        </div>
      </div>
    </div>
    <div class="ui-pagination-total-pages">всего: {{ totalPages_data }}</div>
  </div>
</template>

<script>
export default {
  name: "UIPagination",
  data(){
    return{
      isPagination: false,
      currentPage_data: 1,
      totalPages_data: 1,
      pageInputTimeoutID: undefined,
      inputPage_data: 1,
    }
  },
  props: {
    active: Boolean,
    totalPages: Number,
    currentPage: Number,
  },
  watch: {
    totalPages: function (){
      this.isPagination = this.totalPages > 1
      this.totalPages_data = this.totalPages
    },
    currentPage: function (){
      this.currentPage_data = this.currentPage
      this.inputPage_data = this.currentPage
    },
  },
  methods: {
    onPageInput() {
      if (this.pageInputTimeoutID) {
        clearTimeout(this.pageInputTimeoutID)
      }
      this.pageInputTimeoutID = setTimeout(() => {
        const bad = isNaN(this.inputPage_data) || this.inputPage_data > this.totalPages || this.inputPage_data < 1
        if (bad) {
          return 0
        }
        this.currentPage_data = this.inputPage_data
        this.$emit('page-changed', this.inputPage_data)
      }, 1000)
    },
    onNextButton(){
      this.$emit('page-changed', this.currentPage_data + 1)
    },
    onPrevButton(){
      this.$emit('page-changed', this.currentPage_data - 1)
    },
  },
}
</script>

<style scoped>
.ui-pagination {
  background-color: var(--color-level-1);
  height: 82px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}

.ui-pagination-paginator {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: row;
}

.ui-pagination-next-page,
.ui-pagination-prev-page {
  width: 25%;
}

.ui-pagination-next-page-butt,
.ui-pagination-prev-page-butt {
  cursor: pointer;
  width: 100%;
  height: 100%;
}

.ui-pagination-next-page-butt:hover,
.ui-pagination-prev-page-butt:hover {
  background-color: var(--color-hover);
}

.ui-pagination-next-page-butt,
.ui-pagination-prev-page-butt,
.ui-pagination-pages-input {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.ui-pagination-pages-input {
  width: 50%;
  background-color: var(--color-level-2);
}

.ui-pagination-pages-input > input {
  border: none;
  background-color: transparent;
  width: inherit;
  height: inherit;
  text-align: center;
  font-size: 1.2rem;
  box-sizing: border-box;
}

.ui-pagination-total-pages{
  padding-bottom: 6px;
}
</style>