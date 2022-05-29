<script lang="ts">
	import { browser } from '$app/env';
	import Overlay from '$lib/components/overlay.svelte';
	import { onDestroy, onMount } from 'svelte';

	let isActive = false;
	let titled = '';
	let question = '';
	let resolver: (value: boolean | PromiseLike<boolean>) => void;

	async function plugin(title: string, warningText?: string): Promise<boolean> {
		titled = title;
		question = warningText || 'Are you sure?';
		isActive = true;
		return new Promise((resolve) => {
			resolver = resolve;
		});
	}

	onMount(() => {
		window.$confirm = plugin;
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		window.$confirm = undefined;
		resolveResolver(false);
	});

	function resolveResolver(state: boolean) {
		isActive = false;
		if (resolver) {
			resolver(state);
		}
	}

	function onClose() {
		resolveResolver(false);
	}

	function onYes() {
		resolveResolver(true);
	}
</script>

{#if isActive}
	<Overlay {onClose}>
		<div class="confirm">
			<div class="attention">
				<b class="title">{titled}</b>
				<div class="question">{question}</div>
			</div>
			<div class="accept">
				<div class="no" on:click={onClose}>no</div>

				<span class="divider" />

				<div class="yes" on:click={onYes}>yes</div>
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
			.no,
			.yes {
				&:hover {
					background-color: var(--color-hover);
				}
				cursor: pointer;
			}
		}
	}
</style>
