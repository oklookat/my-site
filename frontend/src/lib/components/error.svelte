<script lang="ts">
	import { onMount } from 'svelte';

	export let imgURL: string | undefined = undefined;
	export let code: string;
	export let short: string;

	onMount(() => {
		if (imgURL) {
			imgURL = '/ineed/error/' + imgURL;
		}
	});
</script>

<div class="error">
	{#if imgURL}
		<img alt="500" src={imgURL} />
	{/if}
	<div class="sign rainbow-text">{code}</div>
	<div class="short" style="font-weight: bold;">({short})</div>
	<div class="long">
		<slot />
	</div>
</div>

<style lang="scss">
	.error {
		width: 100%;
		max-width: 720px;
		display: flex;
		flex-direction: column;
		gap: 12px;
		img {
			padding-top: 12px;
			max-width: 100%;
		}
		img,
		.sign,
		.short {
			align-self: center;
			justify-self: center;
		}
		.sign {
			font-size: 2rem;
		}
		:global(a) {
			color: aqua;
		}
		.long {
			padding-top: 12px;
			width: 100%;
			height: 100%;
			display: flex;
			flex-direction: column;
			gap: 12px;
		}
	}
</style>
