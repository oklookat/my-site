<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    // player
    import ElvenPlayer from "@/plugins/elvenPlayer";
    import Player from "@/plugins/elvenPlayer/index.svelte";
    // choose
    import { ElvenChoose } from "@/plugins/elvenChoose";

    // player
    let player: ElvenPlayer;
    let playerReady = false;
    let playerActive = false;
    // choose
    new ElvenChoose();

    function startPlugins() {
        // player
        player = new ElvenPlayer();
        playerReady = true;
    }

    function destroyPlugins() {
        // player
        playerReady = false;
        player.destroy();
        player = null;
    }

    onMount(() => {
        startPlugins();
    });

    onDestroy(() => {
        destroyPlugins();
    });
</script>

<div class="service">
    {#if playerReady}
        <Player core={player} onActiveChanged={(val) => (playerActive = val)} />
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
    }
</style>
