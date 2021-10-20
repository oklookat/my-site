<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import AuthAdapter from "@/common/adapters/Main/AuthAdapter";

  onMount(async () => {
    AuthAdapter.logout();
    push("/login");
  });
</script>

<div class="logout" />

<style>
</style>
