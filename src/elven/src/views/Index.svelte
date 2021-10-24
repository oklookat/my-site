<script lang="ts">
    import UserAdapter from "@/common/adapters/Main/UserAdapter";
    import type { IUser } from "@/types/user";
    import { onMount } from "svelte";

    let userData: IUser;
    let userDataLoaded = false;
    onMount(() => {
        UserAdapter.getMe().then((user) => {
            userData = user;
            userDataLoaded = true;
        });
    });
</script>

<div class="elven-hello">
    {#if userDataLoaded}
        <div class="account__container">
            <a href="#/account">{userData.username}</a>
        </div>
    {/if}
    <!-- <RouterLink class="header-mobile-menu-item" :to="{ name: 'Settings' }">Настройки</RouterLink>
        <RouterLink class="header-mobile-menu-item" :to="{ name: 'Logout' }">Выход</RouterLink> -->
</div>

<style>
    .account__container {
        background-color: var(--color-level-1);
        border-radius: 14px;
        width: 95%;
        height: 100%;
        max-width: 512px;
        margin: auto;
        display: flex;
        flex-direction: column;
        gap: 14px;
    }
</style>
