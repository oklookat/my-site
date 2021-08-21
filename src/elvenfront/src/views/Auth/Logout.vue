<template>
  <div class="logout">
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import AuthAdapter from "@/common/adapters/Main/AuthAdapter";

export default defineComponent({
  name: 'Logout',
  async mounted() {
    await AuthAdapter.logout()
  }
})
</script>

<style scoped>
.logout{
}

</style>