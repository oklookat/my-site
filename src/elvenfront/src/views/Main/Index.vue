<template>
    <!-- <div class="elven-hello">
        <div class="account">
            здесь будет блок с админ-аккаунтом, по клике на который появятся настройки, смена пароля, выход и так далее
        </div>
        <RouterLink class="header-mobile-menu-item" :to="{ name: 'Settings' }">Настройки</RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{ name: 'Logout' }">Выход</RouterLink>
    </div> -->

    <Scroller></Scroller>
</template>

<script setup lang="ts">
import Scroller from "@/common/plugins/ElvenScroller/new/Scroller.vue"
</script>

<style scoped>
</style>