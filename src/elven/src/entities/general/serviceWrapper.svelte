<script lang="ts">
    import ElvenPlayer from "@/plugins/elvenPlayer"
    import Player from "@/plugins/elvenPlayer/index.svelte"
    import { onDestroy, onMount } from "svelte";

    let player: ElvenPlayer;
    let playerReady = false;

    function startPlugins() {
        player = new ElvenPlayer();
        playerReady = true;
    }

    function destroyPlugins() {
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
