<script lang="ts">
	import { browser } from '$app/env';
	import Overlay from '$lib/components/overlay.svelte';
	import { toggleBodyScroll } from '$elven/tools';
	import { onDestroy, onMount } from 'svelte';

	let isActive = false;
	let titled = '';
	let question = '';
	let resolver: (value: boolean | PromiseLike<boolean>) => void;
	let resetNoScroll: () => void;

	async function plugin(title: string, warningText?: string): Promise<boolean> {
		titled = title;
		question = warningText || 'Are you sure?';
		isActive = true;
		return new Promise((resolve) => {
			resolver = resolve;
		});
	}

	onMount(() => {
		resetNoScroll = toggleBodyScroll();
		window.$confirm = plugin;
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		deactivate();
		window.$confirm = undefined;
	});

	function deactivate() {
		resetNoScroll();
		isActive = false;
	}

	function onClickContainer(e: MouseEvent) {
		if (!resolver) {
			return;
		}
		deactivate();
		resolver(false);
	}

	function onClickYes(e: MouseEvent) {
		if (!resolver) {
			return;
		}
		deactivate();
		resolver(true);
	}

	function onClickNo(e: MouseEvent) {
		if (!resolver) {
			return;
		}
		deactivate();
		resolver(false);
	}
</script>

{#if isActive}
	<Overlay onClose={onClickContainer}>
		<div class="confirm">
			<div class="attention">
				<b class="title">{titled}</b>
				<div class="question">{question}</div>
			</div>
			<div class="accept">
				<div class="no" on:click={onClickNo}>no</div>

				<span class="divider" />

				<div class="yes" on:click={onClickYes}>yes</div>
			</div>
		</div>
	</Overlay>
{/if}

<style lang="scss">
	.confirm {
		width: 264px;
		max-width: 90%;
		min-height: 194px;
		background-color: var(--color-level-1);
		border: 1px solid var(--color-border);
		border-radius: var(--border-radius);
		align-self: center;
		justify-self: center;
		display: grid;
		grid-template-columns: 1fr;
		grid-template-rows: 1fr max-content;
		.attention {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 18px;
			padding: 14px;
		}

		.accept {
			height: 52px;
			bottom: 0;
			width: 100%;
			display: flex;
			span {
				height: 100%;
				width: 1px;
				background-color: var(--color-text);
				opacity: 50%;
			}
			div {
				width: 100%;
				height: 100%;
				display: flex;
				align-items: center;
				justify-content: center;
			}
			.no, .yes {
				&:hover {
					background-color: var(--color-hover);
				}
				cursor: pointer;
			}
		}
	}
</style>
