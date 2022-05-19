<script lang="ts">
	import { toggleBodyScroll } from '$lib/tools';

	import { onDestroy, onMount } from 'svelte';

	export let onClose: () => void;

	let setDefScroll: () => void;
	onMount(() => {
		setDefScroll = toggleBodyScroll();
	});

	onDestroy(() => {
		setDefScroll();
	});
</script>

<div class="overlay base__overlay" on:click|self={onClose}>
	<div class="main">
		<div>
			<slot />
		</div>
	</div>
</div>

<style lang="scss">
	.overlay {
		.main {
			user-select: none;
			background-color: var(--color-level-1);
			border-radius: var(--border-radius);
			min-height: 244px;
			width: 224px;
			@media screen and(max-width: 644px) {
				width: 75%;
			}
			> div {
				// item
				:global(*) {
					cursor: pointer;
					border-bottom: 1px solid var(--color-border);
					min-height: 54px;
					width: 100%;
					display: flex;
					flex-direction: column;
					align-items: center;
					justify-content: center;
					&:hover {
						background-color: var(--color-hover);
					}
				}
				width: 100%;
				height: 100%;
			}
		}
	}
</style>
