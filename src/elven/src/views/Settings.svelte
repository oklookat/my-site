<script lang="ts">
    import AccountSettings from "@/components/AccountSettings.svelte";
    import ToolbarBig from "@/ui/ToolbarBig.svelte";
</script>

<div class="settings">
    <ToolbarBig>
        <a class="block toolbar__logout" href="#/logout">logout</a>
    </ToolbarBig>

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
</style>
