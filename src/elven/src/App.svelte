<script lang="ts">
  // main style
  import "./assets/global.scss";
  // routing
  import Router, { location } from "svelte-spa-router";
  import routes from "@/routes";
  // components
  import Header from "@/components/Header.svelte";
  import ServiceWrapper2 from "@/components/ServiceWrapper2.svelte";
  // plugins
  import Progress from "@/plugins/ElvenProgress/Progress.svelte";
  import Notify from "@/plugins/ElvenNotify/Notify.svelte";

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
    <div style="height: 16px;" />
    <Router {routes} />
    <div style="height: 64px;" />
  </div>
  <Notify />
  {#if isNotAuth}
    <div class="service">
      <ServiceWrapper2 />
    </div>
  {/if}
</div>

<style lang="scss">
  .container {
    min-height: 100vh;
    word-break: break-word;
    display: flex;
    flex-direction: column;
  }

  .content {
    height: 100%;
    width: 100%;
    font-size: 1.1rem;
    line-height: 1.46rem;
    letter-spacing: 0.0007rem;
  }

  .service {
    width: 100%;
    z-index: 7777;
  }

  .service {
    bottom: 0;
  }
</style>
