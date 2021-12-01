<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import AuthAdapter from "@/adapters/AuthAdapter";

  onMount(async () => {
    await AuthAdapter.logout();
    push("/login");
  });
</script>

<div class="logout" />