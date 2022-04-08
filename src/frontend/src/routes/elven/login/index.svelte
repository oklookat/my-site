<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';
	import { goto } from '$app/navigation';

	let username = '';
	let password = '';
	let loginButton: HTMLButtonElement;

	let loginErr: null | string = null;

	$: onUP(username), onUP(password);
	function onUP(v) {
		if (loginErr) {
			loginErr = null;
		}
	}

	async function makeLogin() {
		const jsond = JSON.stringify({ username: username, password: password });
		try {
			const resp = await fetch('/elven/login', {
				method: 'POST',
				body: jsond
			});
			const respJson = await resp.json();
			loginErr = respJson.loginErr;
		} catch (err) {
			loginErr = err;
		}
		if (!loginErr) {
			await goto('/elven');
		}
	}

	onMount(() => {
		document.addEventListener('keydown', onEnter);
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		document.removeEventListener('keydown', onEnter);
	});

	function onEnter(event: KeyboardEvent) {
		// remove double-login when focused on 'log in' button (pressed enter by document event + pressed enter on log in)
		if (event.target === loginButton) {
			return;
		}
		const byEnterAndValid = event.code === 'Enter' && username.length > 0 && password.length > 0;
		if (byEnterAndValid) {
			makeLogin();
		}
	}
</script>

<svelte:head>
	<title>elven: login</title>
</svelte:head>

<div class="login">
	<div class="login__inputs">
		<div class="login__logo logo__text">elven</div>
		<input type="text" name="username" placeholder="username" bind:value={username} />
		<input type="password" name="password" placeholder="password" bind:value={password} />
	</div>
	<button
		disabled={!username || !password}
		class="login__send"
		bind:this={loginButton}
		on:click={makeLogin}
	>
		ok
	</button>
	{#if loginErr}
		{loginErr}
	{/if}
</div>

<style lang="scss">
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
			> input {
				height: 48px;
				width: 264px;
				box-shadow: 0 0 19px 0 rgba(34, 60, 80, 0.02);
			}
		}
		&__send {
			border: 1px solid rgba(0, 0, 0, 0.1);
			box-shadow: 0 0 19px 0 rgba(34, 60, 80, 0.02);
			background-color: var(--color-level-1);
			font-size: 1.2rem;
			width: 264px;
			height: 48px;
		}
	}
</style>
