<script lang="ts">
	import { browser } from '$app/env';

	import { randomColors } from '$elven/tools';
	import { onMount } from 'svelte';

	export let name: string;
	export let path: string;

	let toyEL: HTMLAnchorElement;
	onMount(() => {
		if (!browser) {
			return;
		}
		const color = randomColors();
		toyEL.style.color = color.text;
		toyEL.style.backgroundColor = color.background;
		toyEL.style.border = `2px solid ${randomColors().background}`;
		toyEL.style.display = 'flex';
	});
</script>

<a class="toy" href={path} bind:this={toyEL}>
	<div class="name">
		{name}
	</div>
	<div class="description">
		<slot>Description</slot>
	</div>
</a>

<style lang="scss">
	.toy {
		cursor: pointer;
		border-radius: 4px;
		background-color: var(--color-level-1);
		min-height: 64px;
		width: 100%;
		padding: 8px;
		display: none;
		flex-direction: column;
		gap: 14px;
		.name {
			font-weight: bold;
			font-size: 1.8rem;
		}
	}
</style>
