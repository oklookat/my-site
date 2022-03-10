<script lang="ts">
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import NetworkAuth from "../../network/network_auth";

  onMount(async () => {
    await NetworkAuth.logout();
    push("/login");
  });
</script>

<div class="logout" />
