<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	//
	import NetworkAuth from '$lib/network/network_auth';

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
			window.$notify.add({ message: 'Failed to logout. Try later.' });
			await goto('/elven');
			return;
		}
		await goto('/elven/login');
	});
</script>

<svelte:head>
	<title>elven: logout</title>
</svelte:head>
