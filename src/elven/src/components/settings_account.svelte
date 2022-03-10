<script lang="ts">
    import UserAdapter from "@/network/network_user";
    import Overlay from "@/components/overlay.svelte";
    import type { User, UserChange } from "@/types/user";
    import { usernameValidate, passwordValidate } from "@/types/user";
    import { onMount } from "svelte";

    let userData: User;
    let userDataLoaded = false;

    let changeOverlayActive = false;
    let newValueValid = false;
    let passwordConfirmValid = false;
    let changer: UserChange = {
        what: "username",
        password: "",
        newValue: "",
    };

    onMount(() => {
        getMe();
    });

    async function getMe() {
        try {
            const user = await UserAdapter.getMe();
            userData = user;
            userDataLoaded = true;
        } catch (err) {
            userDataLoaded = false;
        }
    }

    function changeReset() {
        changer.newValue = "";
        changer.password = "";
    }

    function activeUsername() {
        changeOverlayActive = true;
        changer.what = "username";
        changeReset();
    }

    function activePassword() {
        changeOverlayActive = true;
        changer.what = "password";
        changeReset();
    }

    // on username change input
    function onUsernameInput() {
        const username = changer.newValue;
        newValueValid = usernameValidate(username);
    }

    // on password change input
    function onPasswordInput() {
        const password = changer.newValue;
        newValueValid = passwordValidate(password);
    }

    // on password confirm input
    function onPasswordConfirmInput() {
        const password = changer.password;
        passwordConfirmValid = passwordValidate(password);
    }

    // hook before changing username
    function beforeUsername(): boolean {
        if (changer.newValue === userData.username) {
            const message = "You already have this username.";
            window.$notify.add({ message });
            return true;
        }
        return false;
    }

    // hook before changing password
    function beforePassword(): boolean {
        return false;
    }

    async function change() {
        if (!(newValueValid && passwordConfirmValid)) {
            return;
        }
        let err: boolean;
        switch (changer.what) {
            case "username":
                err = beforeUsername();
                break;
            case "password":
                err = beforePassword();
                break;
        }
        if (err) {
            return;
        }
        await UserAdapter.change(changer);
        await getMe();
    }
</script>

{#if userDataLoaded}
    <div class="block account">
        <div class="big">{userData.username}</div>
        <div style="cursor: pointer;" on:click={() => activeUsername()}>
            change username
        </div>
        <div style="cursor: pointer;" on:click={() => activePassword()}>
            change password
        </div>
    </div>
{/if}

{#if changeOverlayActive}
    <Overlay onClose={() => (changeOverlayActive = false)}>
        <div class="overlay">
            <div class="overlay__change">
                <input
                    type="password"
                    placeholder="password"
                    bind:value={changer.password}
                    on:input={onPasswordConfirmInput}
                />
                {#if changer.what === "username"}
                    <input
                        type="text"
                        placeholder="new username"
                        bind:value={changer.newValue}
                        on:input={onUsernameInput}
                    />
                {:else if changer.what === "password"}
                    <input
                        type="password"
                        placeholder="new password"
                        bind:value={changer.newValue}
                        on:input={onPasswordInput}
                    />
                {/if}
                <button
                    disabled={!(newValueValid && passwordConfirmValid)}
                    on:click={change}>change</button
                >
            </div>
        </div>
    </Overlay>
{/if}

<style lang="scss">
    .account {
        height: 122px;
    }

    .overlay {
        &__change {
            width: 100%;
            height: 100%;
            margin-top: 24px;
            display: flex;
            flex-direction: column;
            gap: 16px;
            > input {
                align-self: center;
                height: 46px;
                max-width: 164px;
            }
            > button {
                margin-top: 8px;
                height: 36px;
                width: 84px;
                align-self: center;
            }
        }
    }
</style>
