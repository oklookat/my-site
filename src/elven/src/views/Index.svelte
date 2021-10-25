<script lang="ts">
</script>

<div class="elven-hello">
    <a href="#/settings" class="settings__link">settings</a>
    <!-- <RouterLink class="header-mobile-menu-item" :to="{ name: 'Settings' }">Настройки</RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{ name: 'Logout' }">Выход</RouterLink> -->
</div>

<style>
    .settings__link {
        background-color: var(--color-level-1);
        border-radius: 14px;
        width: 95%;
        min-height: 64px;
        max-width: 512px;
        margin: auto;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 14px;
    }
</style>
