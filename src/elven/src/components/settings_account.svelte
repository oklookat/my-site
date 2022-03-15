<script lang="ts">
    import UserAdapter from "@/network/network_user";
    import type { User, UserChange } from "@/types/user";
    import { usernameValidate, passwordValidate } from "@/types/user";
    import { onMount } from "svelte";

    onMount(() => {
        getMe();
    });

    /** is user data loaded? */
    let isUserDataLoaded = false;

    /** current user */
    let user: User;

    /** get current user data */
    async function getMe() {
        try {
            user = await UserAdapter.getMe();
            isUserDataLoaded = true;
        } catch (err) {
            isUserDataLoaded = false;
        }
    }

    /** change credentials active */
    let isChangeCredentials = false;

    /** changes data */
    let changer: UserChange = {
        what: "username",
        password: "",
        newValue: "",
    };

    /** set what we changing */
    function setChanger(what: "username" | "password") {
        isChangeCredentials = !isChangeCredentials;
        if (!isChangeCredentials) {
            return;
        }
        changer.what = what;
        changer.newValue = "";
        changer.password = "";
    }

    let newValueValid = false;

    /** on new username input */
    function onNewUsernameInput() {
        const username = changer.newValue;
        newValueValid = usernameValidate(username);
    }

    let passwordConfirmValid = false;

    /** on password confirm input */
    function onPasswordInput() {
        const password = changer.password;
        passwordConfirmValid = passwordValidate(password);
    }

    /** on new password  input */
    function onNewPasswordInput() {
        const password = changer.newValue;
        newValueValid = passwordValidate(password);
    }

    /** change username or password depending on changer values */
    async function changeCredentials() {
        if (!(newValueValid && passwordConfirmValid)) {
            return;
        }
        let isError: boolean;
        switch (changer.what) {
            case "username":
                const isAlreadyHaveThisUsername =
                    changer.newValue === user.username;
                isError = isAlreadyHaveThisUsername;
                break;
            case "password":
                isError = false;
                break;
        }
        if (isError) {
            return;
        }
        await UserAdapter.change(changer);
        await getMe();
    }
</script>

{#if isUserDataLoaded}
    <div class="block account">
        <div class="big">{user.username}</div>
        <div class="pointer" on:click={() => setChanger("username")}>
            change username
        </div>
        <div class="pointer" on:click={() => setChanger("password")}>
            change password
        </div>
        {#if isChangeCredentials}
            <div class="change-credentials">
                <input
                    type="password"
                    placeholder="password"
                    bind:value={changer.password}
                    on:input={onPasswordInput}
                />
                {#if changer.what === "username"}
                    <input
                        type="text"
                        placeholder="new username"
                        bind:value={changer.newValue}
                        on:input={onNewUsernameInput}
                    />
                {:else if changer.what === "password"}
                    <input
                        type="password"
                        placeholder="new password"
                        bind:value={changer.newValue}
                        on:input={onNewPasswordInput}
                    />
                {/if}
                <div
                    class="submit button"
                    disabled={!(newValueValid && passwordConfirmValid)}
                    on:click={changeCredentials}
                >
                    change
                </div>
            </div>
        {/if}
    </div>
{/if}

<style lang="scss">
    .account {
        padding: 14px;
        height: max-content;
        .change-credentials {
            width: 100%;
            height: 100%;
            margin-top: 24px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            gap: 16px;
            > input {
                height: 46px;
                max-width: 164px;
            }
            > .submit {
                height: 42px;
                width: 84px;
            }
        }
    }
</style>
