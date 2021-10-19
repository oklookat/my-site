<script lang="ts">
    import { Route } from "svelte-router-spa";
    import Header from "@/components/parts/Header.svelte";
    export let currentRoute;
    const params = {};
</script>

<main>
    <div class="container">
        <elven-notify />
        <div class="service">
            <Header />
            <elven-progress />
        </div>
        <div class="content">
            <div style="height: 16px;" />
            <Route {currentRoute} {params} />
            <div style="height: 64px;" />
        </div>
        <div class="service-2">
            <elven-player />
        </div>
    </div>
</main>

<style>
</style>
