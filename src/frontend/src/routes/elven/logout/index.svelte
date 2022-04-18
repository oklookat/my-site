<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	//
	import NetworkAuth from '$lib_elven/network/network_auth';
	import Utils from '$lib_elven/tools';

	onMount(async () => {
		let isError = false;
		try {
			const resp = await NetworkAuth.logout();
			// if http error and error not 'user not authorized'
			isError = !resp.ok && resp.status !== 401;
		} catch (err) {
			isError = true;
		}
		if (isError) {
			await goto('/elven');
			return;
		}
		await goto('/elven/login');
	});
</script>

<svelte:head>
	<title>{Utils.setTitleElven('logout')}</title>
</svelte:head>
