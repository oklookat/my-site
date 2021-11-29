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
  //import Notify from "@/plugins/ElvenNotify/Notify.svelte";
  import Notify from "@/plugins/ElvenNotify/Notify.svelte";

  let isNotAuth = $location !== "/login" && $location !== "/logout";
  location.subscribe((value) => {
    isNotAuth = value !== "/login" && value !== "/logout";
  });
  // TODO: check user not by local storage. On secured routes get user by request to something like /users/me and check isAdmin field
</script>

<div class="container">
  <div id="progress">
    <Progress />
  </div>
  {#if isNotAuth}
    <div class="header">
      <Header />
    </div>
  {/if}
  <div class="content">
    <div style="height: 16px;" />
    <Router {routes} />
    <div style="height: 64px;" />
  </div>
  <Notify />
  {#if isNotAuth}
    <div class="service-2">
      <ServiceWrapper2 />
    </div>
  {/if}
</div>
