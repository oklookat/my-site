<script lang="ts">
    import ElvenPlayerPlugin from "@/plugins/ElvenPlayer/ElvenPlayerPlugin";
    import { onDestroy, onMount } from "svelte";

    let player: ElvenPlayerPlugin;
    onMount(() => {
        player = new ElvenPlayerPlugin();
    });
    onDestroy(() => {
        player.destroy();
    });
</script>

<div id="elven__player"></div>
