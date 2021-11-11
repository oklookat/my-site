<script lang="ts">
  import "./assets/global.scss";
  import { onDestroy, onMount } from "svelte";
  import Router from "svelte-spa-router";
  import { location } from "svelte-spa-router";
  import routes from "./routes";
  import ElvenNotifyPlugin from "@/plugins/ElvenNotify/ElvenNotifyPlugin";
  import Header from "@/components/parts/Header.svelte";
  import ServiceWrapper2 from "./components/parts/ServiceWrapper2.svelte";
  import Progress from "./plugins/ElvenProgress/Progress.svelte";

  let isNotAuth = $location !== "/login" && $location !== "/logout";
  location.subscribe((value) => {
    isNotAuth = value !== "/login" && value !== "/logout";
  });

  // TODO: rewrite plugins like Progress
  // init plugins
  let elvenNotifyEL;
  let elvenNotify: ElvenNotifyPlugin;

  onMount(() => {
    elvenNotify = new ElvenNotifyPlugin(elvenNotifyEL);
  });

  onDestroy(() => {
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
  <div id="elven__progress">
    <Progress />
  </div>
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
