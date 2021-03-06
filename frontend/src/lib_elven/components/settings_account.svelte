<script lang="ts">
	import { onMount } from 'svelte';
	import NetworkUser from '$elven/network/user';
	import type { User, UserChange } from '$elven/types/user';
	import { usernameValidate, passwordValidate } from '$elven/types/user';

	const networkUser = new NetworkUser('');

	/** is user data loaded? */
	let isUserDataLoaded = false;

	/** current user */
	let user: User;

	/** change credentials active */
	let isChangeCredentials = false;

	/** changes data */
	let changer: UserChange = {
		what: 'username',
		password: '',
		newValue: ''
	};

	onMount(async () => {
		await getMe();
	});

	/** get current user data */
	async function getMe() {
		try {
			const resp = await networkUser.getMe();
			if (resp.ok) {
				user = await resp.json();
				isUserDataLoaded = true;
			}
		} catch (err) {
			isUserDataLoaded = false;
		}
	}

	/** set what we changing */
	function setChanger(what: 'username' | 'password') {
		isChangeCredentials = !isChangeCredentials;
		if (!isChangeCredentials) {
			return;
		}
		changer.what = what;
		changer.newValue = '';
		changer.password = '';
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
		let isHasError = false;
		switch (changer.what) {
			case 'username':
				const isAlreadyHaveThisUsername = changer.newValue === user.username;
				isHasError = isAlreadyHaveThisUsername;
				break;
			case 'password':
				isHasError = false;
				break;
		}
		if (isHasError) {
			return;
		}
		await NetworkUser.change(changer);
		await getMe();
	}
</script>

{#if isUserDataLoaded}
	<div class="account">
		<b class="username">{user.username}</b>
		<div style="cursor: pointer" on:click={() => setChanger('username')}>change username</div>
		<div style="cursor: pointer" on:click={() => setChanger('password')}>change password</div>

		{#if isChangeCredentials}
			<div class="change-credentials">
				<input
					type="password"
					placeholder="password"
					bind:value={changer.password}
					on:input={onPasswordInput}
				/>
				{#if changer.what === 'username'}
					<input
						type="text"
						placeholder="new username"
						bind:value={changer.newValue}
						on:input={onNewUsernameInput}
					/>
				{:else if changer.what === 'password'}
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
					ok
				</div>
			</div>
		{/if}
	</div>
{/if}

<style lang="scss">
	.account {
		background-color: var(--color-level-1);
		font-size: 1.3rem;
		padding: 14px;
		.username {
			font-size: 1.4rem;
		}
	}

	.change-credentials {
		width: 100%;
		height: 100%;
		padding-top: 24px;
		display: flex;
		flex-direction: column;
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
</style>
