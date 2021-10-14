<template>
  <div class="ui-pagination" v-if="active">
    <div class="ui-pagination__paginator">
      <div class="ui-pagination__prev-page">
        <div
          class="ui-pagination__prev-page-butt"
          v-if="currentPageData !== 1"
          v-on:click="onPrevButton"
        >назад</div>
      </div>
      <div class="ui-pagination__pages-input">
        <input
          class="ui-pagination__pages-input-num"
          type="text"
          placeholder="Номер страницы"
          v-model="inputPage"
          @input="onPageInput"
        />
      </div>
      <div class="ui-pagination__next-page">
        <div
          class="ui-pagination__next-page-butt"
          v-if="currentPageData < totalPagesData"
          v-on:click="onNextButton"
        >вперед</div>
      </div>
    </div>
    <div class="ui-pagination__total-pages">всего: {{ totalPagesData }}</div>
  </div>
</template>

<script setup lang="ts">
import { watch, ref, Ref } from '@vue/runtime-core'

const props = defineProps<{
  totalPages: number,
  currentPage: number
}>()

const emit = defineEmits<{
  (e: 'changed', page: number): void
}>()

const active = ref(false)
const totalPagesData = ref(1)
const currentPageData = ref(1)
const pageInputTimeoutID: Ref<ReturnType<typeof setTimeout> | null> = ref(null)
const inputPage: Ref<number> = ref(1)


watch(() => props.totalPages, (newValue, prevValue) => {
  active.value = newValue > 1
  totalPagesData.value = newValue
})

watch(() => props.currentPage, (newValue, prevValue) => {
  currentPageData.value = newValue
  inputPage.value = newValue
})

function onPageInput() {
  if (pageInputTimeoutID.value) {
    clearTimeout(pageInputTimeoutID.value)
  }
  pageInputTimeoutID.value = setTimeout(() => {
    const bad = isNaN(inputPage.value) || inputPage.value > totalPagesData.value || inputPage.value < 1
    if (bad) {
      return 0
    }
    currentPageData.value = inputPage.value
    emit('changed', inputPage.value)
  }, 1000)
}

function onNextButton() {
  emit('changed', currentPageData.value + 1)
}

function onPrevButton() {
  emit('changed', currentPageData.value - 1)
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

.ui-pagination__paginator {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: row;
}

.ui-pagination__next-page,
.ui-pagination__prev-page {
  width: 25%;
}

.ui-pagination__next-page-butt,
.ui-pagination__prev-page-butt {
  cursor: pointer;
  width: 100%;
  height: 100%;
}

.ui-pagination__next-page-butt:hover,
.ui-pagination__prev-page-butt:hover {
  background-color: var(--color-hover);
}

.ui-pagination__next-page-butt,
.ui-pagination__prev-page-butt,
.ui-pagination__pages-input {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.ui-pagination__pages-input {
  width: 50%;
  background-color: var(--color-level-2);
}

.ui-pagination__pages-input > input {
  border: none;
  background-color: var(--color-hover);
  width: 100%;
  height: inherit;
  text-align: center;
  font-size: 1.2rem;
  box-sizing: border-box;
}

.ui-pagination__total-pages {
  padding-bottom: 6px;
}
</style>