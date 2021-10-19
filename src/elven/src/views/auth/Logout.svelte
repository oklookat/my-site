<main>
  <div class="logout">
  </div>
</main>

<script lang="ts">
import AuthAdapter from "@/common/adapters/Main/AuthAdapter"
import { onMount } from "svelte";

onMount(async () => {
  await AuthAdapter.logout()
})
</script>

<style>
</style>