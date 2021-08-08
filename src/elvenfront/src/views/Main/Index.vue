<template>
<div class="elven-hello">
  <div class="account">здесь будет блок с админ-аккаунтом, по клике на который появятся настройки, смена пароля, выход и так далее</div>
  <RouterLink class="header-mobile-menu-item" :to="{name: 'Settings'}">Настройки</RouterLink>
  <RouterLink class="header-mobile-menu-item" :to="{name: 'Logout'}">Выход</RouterLink>
</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Index',
  components: {}
})
</script>

<style scoped>

</style>