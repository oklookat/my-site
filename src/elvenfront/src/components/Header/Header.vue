<template>
  <div class="header-container">

    <div class="header">
      <nav class="header-navigation">
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Index'}">
          elven
        </RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Articles'}">
          записи
        </RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{name: 'Files'}">
          файлы
        </RouterLink>
      </nav>
    </div>


  </div>
</template>

<script>
import {defineComponent} from 'vue'

export default defineComponent({
  name: 'Header',
  data() {
    return {}
  }
})
</script>

<style>
.header-container {
  background-color: var(--color-header);
  color: var(--color-header-text);
  height: var(--header-height);
  width: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.header{
  height: 100%;
  max-width: 975px;
}

.header-navigation {
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: transparent;
}


.header-mobile-menu-item {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  width: 84px;
  cursor: pointer;
}
</style>