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
    // choose
    const choose = new ElvenChoose();

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

{#if playerReady}
    <Player core={player} />
{/if}
