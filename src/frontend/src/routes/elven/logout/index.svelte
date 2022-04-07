<script lang="ts">
  import { onMount } from "svelte";
  import NetworkAuth from "$lib/network/network_auth";
  import { goto } from "$app/navigation"; 

  onMount(async () => {
    await NetworkAuth.logout();
    await goto("/elven/login")
  });
</script>

<svelte:head>
  <title>elven: logout</title>
</svelte:head>

<div class="logout" />
