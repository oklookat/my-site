<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import Notify from "$lib/plugins/elvenNotify/notify.svelte"
    // player
    import Player from "$lib/plugins/elvenPlayer/components/index.svelte";

    let isPlayerReady = false;
    let playerCore;
    //
    onMount(async () => {
        const PlayerCore = await import("$lib/plugins/elvenPlayer")
        const { ElvenChoose } = await import("$lib/plugins/elvenChoose")
        new ElvenChoose()

        playerCore = new PlayerCore.default()
        isPlayerReady = true
    })

    function destroyPlugins() {
        // player
        //playerCore.destroy();
       // playerCore = null;
    }

    onDestroy(() => {
        destroyPlugins();
    });
</script>

<div class="service">
    <Notify />
    {#if isPlayerReady}
        <Player core={playerCore} />
    {/if}
</div>

<style lang="scss">
    .service {
        display: flex;
        height: fit-content;
        background-color: red;
        bottom: 0;
        width: 100%;
        z-index: 7777;
        position: sticky;
        bottom: 0;
    }
</style>
