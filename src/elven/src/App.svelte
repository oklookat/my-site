<script lang="ts">
  import "./assets/global.css";
  import { onDestroy, onMount } from "svelte";
  import Router from "svelte-spa-router";
  import routes from "./routes";
  import ElvenProgress from "@/common/plugins/ElvenProgress/ElvenProgress"
  import ElvenNotify from "@/common/plugins/ElvenNotify/ElvenNotify";
  import { location } from "svelte-spa-router";
  import ServiceWrapper from "./components/parts/ServiceWrapper.svelte";
  import ServiceWrapper2 from "./components/parts/ServiceWrapper2.svelte";

  let isNotAuth = $location !== "/login" && $location !== "/logout";
  location.subscribe((value) => {
    isNotAuth = value !== "/login" && value !== "/logout";
  });

  // init plugins
  let elvenProgress: ElvenProgress;
  let elvenNotify: ElvenNotify;
  onMount(() => {
    elvenProgress = new ElvenProgress();
    elvenNotify = new ElvenNotify();
  });
  onDestroy(() => {
    elvenProgress.destroy();
    elvenNotify.destroy();
  });
</script>

<div class="container">
  <div id="elven__notify" />
  {#if isNotAuth}
    <ServiceWrapper />
  {/if}
  <div id="elven__progress" />
  <div class="content">
    <div style="height: 16px;" />
    <Router {routes} />
    <div style="height: 64px;" />
  </div>
  {#if isNotAuth}
    <ServiceWrapper2 />
  {/if}
</div>
