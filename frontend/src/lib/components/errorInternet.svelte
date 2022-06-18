<script>
	import { onMount } from 'svelte';

	import Error from './error.svelte';

	let isOffline = false;
	let counter = 0;
	onMount(() => {
		isOffline = !navigator.onLine;
	});

	function check() {
		isOffline = !navigator.onLine;
		counter++;
	}
</script>

<Error code="Internet" short="нет интернета">
	{#if isOffline}
		<div on:click={check}><b>А сейчас? (нажмите чтобы проверить)</b></div>
		<div>Проверено столько раз: {counter}</div>
	{/if}
	{#if isOffline && counter > 0}
		<div><b>Нет, сейчас тоже.</b></div>
	{/if}
	{#if !isOffline}
		<div><b>Вроде есть. Попробуйте перейти на главную.</b></div>
	{/if}
</Error>

<style lang="scss">
	div {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 16px;
	}
</style>
