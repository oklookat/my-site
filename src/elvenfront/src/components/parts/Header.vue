<template>
  <div class="header__container">
    <div class="header__universal">
      <nav class="header__navigation">
        <RouterLink
          class="header__item"
          :to="{ name: 'Index' }"
          :class="{ 'header__item-active': path === '/' }"
        >elven</RouterLink>
        <RouterLink
          class="header__item"
          :to="{ name: 'Articles' }"
          :class="{ 'header__item-active': path.startsWith('/articles') }"
        >articles</RouterLink>
        <RouterLink
          class="header__item"
          :to="{ name: 'Files' }"
          :class="{ 'header__item-active': path.startsWith('/files') }"
        >files</RouterLink>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from '@vue/reactivity'
import { useRouter, useRoute } from 'vue-router'
const router = useRouter()
const route = useRoute()
let path = ref('/')
router.beforeEach((to, from, next) => {
  path.value = to.path
  next()
})
</script>

<style>
.header__container {
  /* background-color: var(--color-header); */
  font-weight: bold;
  color: var(--color-header-text);
  height: var(--header-height);
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header__universal {
  height: 100%;
  width: 100%;
  max-width: 975px;
}

.header__item {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 80%;
  width: 84px;
  cursor: pointer;
}

.header__item-active {
  border-radius: 4px;
  background-color: var(--color-header-active);
}

.header__navigation {
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: transparent;
}
</style>