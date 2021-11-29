<script lang="ts">
    import { onMount } from "svelte";
    import AccountSettings from "@/components/AccountSettings.svelte";

    onMount(() => {});
</script>

<div class="settings">
    <div class="toolbar">
        <a class="block toolbar__logout" href="#/logout">Logout</a>
    </div>
    <AccountSettings />
</div>

<style lang="scss">
    .settings {
        width: 95%;
        height: 100%;
        min-height: 64px;
        max-width: 512px;
        margin: auto;
        display: flex;
        flex-direction: column;
        gap: 14px;
    }

    .toolbar {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        grid-template-rows: 1fr;
        &__logout {
            height: 42px;
            display: flex;
            justify-content: center;
            align-items: center;
        }
    }
</style>
