<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import Adapter from "../adapter";

  onMount(async () => {
    await Adapter.logout();
    push("/login");
  });
</script>

<div class="logout" />