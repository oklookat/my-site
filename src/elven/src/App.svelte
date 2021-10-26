<script lang="ts">
  import "./assets/global.css";
  import { onDestroy, onMount } from "svelte";
  import Router from "svelte-spa-router";
  import { location } from "svelte-spa-router";
  import routes from "./routes";
  import ElvenProgressPlugin from "@/plugins/ElvenProgress/ElvenProgressPlugin";
  import ElvenNotifyPlugin from "@/plugins/ElvenNotify/ElvenNotifyPlugin";
  import Header from "@/components/parts/Header.svelte";
  import ServiceWrapper2 from "./components/parts/ServiceWrapper2.svelte";

  let isNotAuth = $location !== "/login" && $location !== "/logout";
  location.subscribe((value) => {
    isNotAuth = value !== "/login" && value !== "/logout";
  });

  // init plugins
  let elvenProgressEL;
  let elvenNotifyEL;
  let elvenProgress: ElvenProgressPlugin;
  let elvenNotify: ElvenNotifyPlugin;

  onMount(() => {
    elvenProgress = new ElvenProgressPlugin(elvenProgressEL);
    elvenNotify = new ElvenNotifyPlugin(elvenNotifyEL);
  });

  onDestroy(() => {
    elvenProgress.destroy();
    elvenNotify.destroy();
  });
  // TODO: check user not by local storage. On secured routes get user by request to something like /users/me and check isAdmin field
</script>

<div class="container">
  {#if isNotAuth}
    <div class="global__header">
      <Header />
    </div>
  {/if}
  <div id="elven__progress" bind:this={elvenProgressEL} />
  <div class="content">
    <div style="height: 16px;" />
    <Router {routes} />
    <div style="height: 64px;" />
  </div>
  <div id="elven__notify" bind:this={elvenNotifyEL} />
  {#if isNotAuth}
    <div class="service-2">
      <ServiceWrapper2 />
    </div>
  {/if}
</div>
