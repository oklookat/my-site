<script lang="ts">
    import AccountSettings from "@/entities/settings/account.svelte";
    import ToolbarBig from "@/ui/toolbarBig.svelte";
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
