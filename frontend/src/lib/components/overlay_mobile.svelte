<script lang="ts">
	import { fly } from 'svelte/transition';
	import Overlay from './overlay.svelte';
	export let title = "Title"
	export let onClose: (e: MouseEvent) => void;
</script>

<Overlay {onClose}>
	<div transition:fly={{ y: 1000, duration: 224 }} class="mobile">
		<div class="title">
			{title}
		</div>
		<div class="content">
			<slot />
		</div>
	</div>
</Overlay>

<style lang="scss">
	.mobile {
		align-self: flex-end;
		border-top-left-radius: 8px;
		border-top-right-radius: 8px;
		background-color: var(--color-body);
		width: 100%;
		min-height: 50%;
		height: max-content;
		display: grid;
		grid-template-columns: 1fr;
		grid-template-rows: max-content 1fr;
		.title {
			padding: 6px;
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
