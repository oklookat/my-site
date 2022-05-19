<script lang="ts">
	import { browser } from '$app/env';
	import { onDestroy, onMount } from 'svelte';
import { t } from '$lib/locale';

	let isActive = false;
	let titled = '';
	let question = '';

	let resolver: (value: boolean | PromiseLike<boolean>) => void;

	let defaultOverflow = '';
	const toggleOverflow = createOverflowToggler();

	async function impl(title: string, warningText?: string): Promise<boolean> {
		titled = title;
		if (warningText) {
			question = warningText;
		} else {
			question = $t('elven.plugins.elvenChoose.areYouSure');
		}
		toggleOverflow();
		isActive = true;
		return new Promise((resolve) => {
			resolver = resolve;
		});
	}

	onMount(() => {
		defaultOverflow = document.body.style.overflow;
		window.$confirm = impl;
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		deactivate();
		window.$confirm = undefined;
	});

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

	function deactivate() {
		document.body.style.overflow = defaultOverflow;
		isActive = false;
	}

	/** create default / hidden overflow setter */
	function createOverflowToggler(): () => void {
		if (!browser) {
			return () => {};
		}

		let isDef = true;

		return () => {
			if (isDef) {
				document.body.style.overflow = 'hidden';
				isDef = false;
				return;
			}
			document.body.style.overflow = defaultOverflow;
			isDef = true;
		};
	}
</script>

{#if isActive}
	<div class="confirm" on:click|stopPropagation|self={onClickContainer}>
		<div class="confirm__second">
			<div class="confirm__title">{titled}</div>
			<div class="confirm__question">{question}</div>
			<div class="confirm__ny">
				<div class="confirm__no" on:click={onClickNo}>{$t('elven.plugins.elvenChoose.no')}</div>

				<div class="confirm__divider" />

				<div class="confirm__yes" on:click={onClickYes}>{$t('elven.plugins.elvenChoose.yes')}</div>
			</div>
		</div>
	</div>
{/if}

<style lang="scss">
	.confirm {
		box-sizing: border-box;
		background-color: rgba(0, 0, 0, 0.4);

		z-index: 9998;
		width: 100%;
		height: 100%;
		position: fixed;

		top: 0;
		right: 0;
		bottom: 0;
		left: 0;

		display: flex;
		align-items: center;
		justify-content: center;

		&__second {
			background-color: var(--color-level-1);
			border-radius: var(--border-radius);

			width: 95%;
			max-width: 244px;

			height: fit-content;
			min-height: 164px;
			max-height: 35%;

			word-break: break-word;
			overflow: hidden;

			display: grid;
			grid-template-columns: 100%;
			grid-template-rows: max-content 1fr auto;
		}
		$padding: 8px;
		&__title {
			border-radius: var(--border-radius) var(--border-radius) 0 0;
			letter-spacing: 0.05rem;
			min-height: 38px;
			padding: $padding;
			font-size: 1.2rem;
			display: flex;
			align-items: center;
			justify-content: center;
		}
		&__question {
			font-size: 1rem;
			line-height: 1.4rem;
			width: 100%;
			max-width: 100%;
			padding: $padding;
			overflow: auto;
			display: flex;
			justify-content: center;
			align-items: center;
		}

		&__ny {
			width: 100%;
			display: flex;
			height: 44px;
			margin-bottom: auto;
		}
		&__divider {
			pointer-events: none;
			height: 100%;
			width: 1px;
			border: 1px solid var(--color-border);
		}
		&__no,
		&__yes {
			font-size: 1.2rem;
			cursor: pointer;
			width: 100%;
			height: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
			&:hover {
				background-color: var(--color-hover);
			}
		}
	}
</style>
