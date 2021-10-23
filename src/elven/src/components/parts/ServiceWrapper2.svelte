<script lang="ts">
    import ElvenPlayerPlugin from "@/common/plugins/ElvenPlayer/ElvenPlayerPlugin";
    import { onDestroy, onMount } from "svelte";

    let player: ElvenPlayerPlugin;
    onMount(() => {
        player = new ElvenPlayerPlugin();
    });
    onDestroy(() => {
        player.destroy();
    });
</script>

<div class="service-2">
    <div id="elven__player"></div>
</div>
