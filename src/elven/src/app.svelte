<script lang="ts">
  // main style
  import "./assets/global.scss";
  // plugins
  import Progress from "@/plugins/elvenProgress/progress.svelte";
  import Notify from "@/plugins/elvenNotify/notify.svelte";
  // routing
  import Router, { location } from "svelte-spa-router";
  import routes from "@/routes";
  // components
  import Header from "@/entities/general/header.svelte";
  import ServiceWrapper from "@/entities/general/serviceWrapper.svelte";

  let isNotAuth = $location !== "/login" && $location !== "/logout";
  location.subscribe((value) => {
    isNotAuth = value !== "/login" && value !== "/logout";
  });
  // TODO: check user not by local storage. On secured routes get user by request to something like /users/me and check isAdmin field
</script>

<div class="container">
  <Progress />
  {#if isNotAuth}
    <Header />
  {/if}
  <div class="content">
    <Router {routes} restoreScrollState={true} />
  </div>
  <Notify />
  {#if isNotAuth}
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
