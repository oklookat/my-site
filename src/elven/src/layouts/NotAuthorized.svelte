<script lang="ts">
  import { Route } from "svelte-router-spa";
  export let currentRoute;
  const params = {};
</script>

<main>
  <div class="container">
    <elven-notify />
    <div class="content">
      <div style="height: 16px;" />
      <Route {currentRoute} {params} />
      <div style="height: 64px;" />
    </div>
  </div>
</main>

<style>
</style>
