<script lang="ts">
	import { cutString } from '$elven/tools';

	import { fly } from 'svelte/transition';
	import Overlay from './overlay.svelte';
	export let title = 'Title';
	export let onClose: (e: MouseEvent) => void;
</script>

<Overlay {onClose}>
	<div transition:fly={{ y: 1000, duration: 224 }} class="mobile">
		<div class="title">
			{cutString(title)}
		</div>
		<div class="content">
			<slot />
		</div>
	</div>
</Overlay>

<style lang="scss">
	.mobile {
		word-break: break-word;
		align-self: flex-end;
		justify-self: center;
		border-top-left-radius: 8px;
		border-top-right-radius: 8px;
		background-color: var(--color-body);
		max-width: 724px;
		width: 100%;
		height: 50%;
		overflow: auto;
		display: grid;
		grid-template-columns: 1fr;
		grid-template-rows: max-content 1fr;
		.title {
			position: sticky;
			top: 0;
			padding: 12px;
			width: 100%;
			min-height: 25px;
			border-radius: inherit;
			border-top: 1px solid var(--color-border);
			background-color: var(--color-level-1);
			font-weight: bold;
			display: flex;
			align-items: center;
			justify-content: center;
		}
	}
</style>
