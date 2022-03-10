<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import { GlobalState } from "@/tools/storage";
    import { push } from "svelte-spa-router";

    onMount(() => {
        GlobalState.isNotFoundPage.set(true);
    });
    onDestroy(() => {
        GlobalState.isNotFoundPage.set(false);
    });

    function goMain() {
        push("#/");
    }
</script>

<div class="notfound" on:click={() => goMain()}>
    <div class="number">404</div>
    <div class="message">page not found</div>
    <div class="go">click anywhere to go to home</div>
</div>

<style lang="scss">
    .notfound {
        width: 100%;
        height: 100vh;
        position: absolute;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 24px;
        .number {
            font-size: 5rem;
            width: auto;
            height: 64px;
            display: flex;
            align-items: center;
            justify-content: center;
        }
    }
</style>
