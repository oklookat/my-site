<script lang="ts">
    import { ElvenNotify } from "@/plugins/ElvenNotify";
    import { onMount } from "svelte";

    let container: HTMLDivElement;
    let plugin: ElvenNotify;
    onMount(() => {
        plugin = new ElvenNotify(container);
    });
</script>

<div class="notify">
    <div class="notify__notifications" bind:this={container} />
</div>

<style lang="scss">
    @mixin _desktop {
        @media (min-width: 765px) {
            @content;
        }
    }

    .notify {
        width: 100%;
        z-index: 9999;
        bottom: 0;
        margin-bottom: 8px;
        position: fixed;
        overflow: hidden;
        @include _desktop() {
            margin-right: 12px;
            height: min-content;
            width: 224px;
            right: 0;
            bottom: 0;
        }
    }
</style>
