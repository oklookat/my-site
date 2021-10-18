<template>
  <div class="logout">
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "@vue/runtime-core"
import AuthAdapter from "@/common/adapters/Main/AuthAdapter"

onMounted(async () => {
  await AuthAdapter.logout()
})
</script>

<style scoped>
</style>