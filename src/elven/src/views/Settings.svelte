<script lang="ts">
    import { onMount } from "svelte";
    import ChangeCredentials from "@/components/parts/ChangeCredentials.svelte";
    import Overlay from "@/components/ui/Overlay.svelte";

    onMount(() => {});
</script>

<div class="settings__container">
    <div class="settings__content">
        <div class="toolbar__one">
            <a class="block account__logout" href="#/logout">Logout</a>
        </div>
        <ChangeCredentials />
    </div>
</div>

<style>
    .settings__container {
        width: 95%;
        height: 100%;
        min-height: 64px;
        max-width: 512px;
        margin: auto;
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    .settings__content {
        display: flex;
        flex-direction: column;
        gap: 14px;
    }

    .toolbar__one {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        grid-template-rows: 1fr;
    }

    .account__logout {
        height: 42px;
        display: flex;
        justify-content: center;
        align-items: center;
    }
</style>
