<script lang="ts">
  import AuthAdapter from "@/common/adapters/Main/AuthAdapter";
  import { onDestroy, onMount } from "svelte";
  import {push} from 'svelte-spa-router'

  let username = "";
  let password = "";
  let loginButton;

  function makeLogin() {
    AuthAdapter.login(username, password)
      .then(() => {
        push("/");
      })
      .catch(() => {
      });
  }

  onMount(() => {
    document.addEventListener("keydown", onEnter);
  });

  onDestroy(() => {
    document.removeEventListener("keydown", onEnter);
  });

  function onEnter(event) {
    // remove double-login when focused on 'log in' button (pressed enter by document event + pressed enter on log in)
    if(event.target === loginButton){
      return
    }
    if (event.code === "Enter" && username.length > 0 && password.length > 0) {
      makeLogin();
    }
  }
</script>

<div class="login-container">
  <div class="login-inputs">
    <div class="logo-text">elven</div>
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
    class="login-button"
    bind:this={loginButton}
    on:click={makeLogin}
  >
    Log in
  </button>
</div>

<style>
  .login-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 18px;
  }

  .login-inputs {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-top: 108px;
    gap: 14px;
  }
  .logo-text {
    margin-bottom: 16px;
  }

  input {
    height: 48px;
    width: 264px;
    box-shadow: 0 0 19px 0 rgba(34, 60, 80, 0.02);
  }

  .login-button {
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
</style>
