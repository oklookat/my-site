<script lang="ts">
  // main style
  import "./assets/global.scss";
  // plugins
  import Progress from "@/plugins/elvenProgress/progress.svelte";
  import Notify from "@/plugins/elvenNotify/notify.svelte";
  // routing
  import Router, { location } from "svelte-spa-router";
  import routes from "@/views/routes";
  // components
  import Header from "@/components/header.svelte";
  import ServiceWrapper from "@/components/service_wrapper.svelte";
  import { GlobalState } from "@/tools/storage";

  // TODO: check user not by local storage.
  // On secured routes get user by request to something like /users/me and check isAdmin field

  // state
  let isAuthPage = false;
  location.subscribe((value) => {
    isAuthPage = value.includes("/login") && value.includes("/logout");
  });
  let is404Page = false;
  GlobalState.isNotFoundPage.subscribe((value) => {
    is404Page = value;
  });
</script>

<div class="container">
  <Progress />
  {#if !isAuthPage && !is404Page}
    <Header />
  {/if}
  <div class="content">
    <Router {routes} restoreScrollState={true} />
  </div>
  <Notify />
  {#if !isAuthPage && !is404Page}
    <ServiceWrapper />
  {/if}
</div>

<style lang="scss">
  .container {
    min-height: 100vh;
    word-break: break-word;
    display: grid;
    grid-template-columns: 1fr;
    // header - content - service
    grid-template-rows: max-content 1fr min-content;
    gap: 16px;
  }

  .content {
    height: 100%;
    width: 100%;
    font-size: 1.1rem;
    line-height: 1.46rem;
    letter-spacing: 0.0007rem;
  }
</style>
