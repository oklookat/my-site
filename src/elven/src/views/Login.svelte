<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import AuthAdapter from "@/adapters/AuthAdapter";

  let username = "";
  let password = "";
  let loginButton;

  function makeLogin() {
    AuthAdapter.login(username, password)
      .then(() => {
        push("/");
      })
      .catch(() => {});
  }

  onMount(() => {
    document.addEventListener("keydown", onEnter);
  });

  onDestroy(() => {
    document.removeEventListener("keydown", onEnter);
  });

  function onEnter(event) {
    // remove double-login when focused on 'log in' button (pressed enter by document event + pressed enter on log in)
    if (event.target === loginButton) {
      return;
    }
    if (event.code === "Enter" && username.length > 0 && password.length > 0) {
      makeLogin();
    }
  }
</script>

<div class="login">
  <div class="login__inputs">
    <div class="login__logo logo__text">elven</div>
    <input
      type="text"
      name="username"
      placeholder="Username"
      bind:value={username}
    />
    <input
      type="password"
      name="password"
      placeholder="Password"
      bind:value={password}
    />
  </div>
  <button
    disabled={!username || !password}
    class="login__send"
    bind:this={loginButton}
    on:click={makeLogin}
  >
    Log in
  </button>
</div>

<style lang="scss">
  input {
    height: 48px;
    width: 264px;
    box-shadow: 0 0 19px 0 rgba(34, 60, 80, 0.02);
  }

  .login {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 18px;
    &__logo {
      margin-bottom: 16px;
    }
    &__inputs {
      display: flex;
      flex-direction: column;
      align-items: center;
      margin-top: 108px;
      gap: 14px;
    }
    &__send {
      border: 1px solid rgba(0, 0, 0, 0.1);
      box-shadow: 0 0 19px 0 rgba(34, 60, 80, 0.02);
      background-color: var(--color-level-1);
      border-radius: 8px;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      user-select: none;
      cursor: pointer;
      font-size: 1.2rem;
      width: 264px;
      height: 48px;
    }
  }
</style>
