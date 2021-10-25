<script lang="ts">
    import UserAdapter from "@/common/adapters/Main/UserAdapter";
    import Overlay from "@/components/ui/Overlay.svelte";
    import type { TUser, TUserChange } from "@/types/user";
    import { usernameValidate, passwordValidate } from "@/types/user";
    import { onMount } from "svelte";
    let userData: TUser;
    let userDataLoaded = false;

    let changeOverlayActive = false;
    let usernameValid = false;
    let passwordValid = false;
    let passwordConfirmValid = false;
    let changer: TUserChange = {
        what: "username",
        password: "",
        newValue: "",
    };

    onMount(() => {
        UserAdapter.getMe().then((user) => {
            userData = user;
            userDataLoaded = true;
        });
    });

    function changeUsername() {
        changeOverlayActive = true;
        changer.what = "username";
    }

    function changePassword() {
        changeOverlayActive = true;
        changer.what = "password";
    }

    function onUsernameInput() {
        const username = changer.newValue;
        usernameValid = usernameValidate(username);
    }

    function onPasswordInput() {
        const password = changer.newValue;
        passwordValid = passwordValidate(password);
    }

    function onPasswordConfirmInput(){
        const password = changer.password;
        passwordValid = passwordValidate(password);
    }
</script>

{#if userDataLoaded}
    <div class="settings__container">
        <div class="settings__content">
            <a class="block account__logout" href="#/logout">Logout</a>
            <div class="block account__settings">
                <div class="big">{userData.username}</div>
                <div style="cursor: pointer;" on:click={() => changeUsername()}>
                    change username
                </div>
                <div style="cursor: pointer;" on:click={() => changePassword()}>
                    change password
                </div>
            </div>
        </div>
    </div>
{/if}

<Overlay
    active={changeOverlayActive}
    on:deactivated={() => (changeOverlayActive = false)}
>
    {#if changer.what === "username"}
        <div class="overlay__username">
            <input
                type="text"
                placeholder="new username"
                bind:value={changer.newValue}
                on:input={onUsernameInput}
            />
            <input
                type="password"
                placeholder="password"
                bind:value={changer.password}
                on:input={onPasswordConfirmInput}
            />
            <button disabled={!(usernameValid && passwordConfirmValid)}
                >confirm</button
            >
        </div>
    {/if}
</Overlay>

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

    .block {
        background-color: var(--color-level-1);
        border-radius: 14px;
        display: flex;
        flex-direction: column;
        justify-content: center;
        gap: 8px;
    }

    .block > * {
        margin-left: 8px;
    }

    .big {
        font-size: 1.5rem;
        font-weight: bold;
        margin-bottom: 6px;
    }

    .account__logout {
        height: 42px;
        display: flex;
        justify-content: center;
        align-items: center;
    }
    .account__settings {
        height: 122px;
    }
    .overlay__username {
        display: flex;
        flex-direction: column;
        gap: 16px;
        height: 100%;
        margin-top: 24px;
    }
    .overlay__username > input {
        height: 46px;
    }
    .overlay__username > button {
        height: 36px;
        width: 84px;
        align-self: center;
    }
</style>
