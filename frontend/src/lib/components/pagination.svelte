<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	

	const dispatch = createEventDispatcher<{
		/** when page changed */
		changed: number;
	}>();

	/** total pages */
	export let total: number;

	/** current page */
	export let current: number;

	/** page from input */
	let inputPage = `${current}`;

	$: watchCurrentPage(current);
	function watchCurrentPage(value: number) {
		if (!value || value < 1) {
			value = 1;
		}
		current = value;
		inputPage = value.toString();
	}

	/** timeout for page input */
	let pageInputTimeoutID: ReturnType<typeof setTimeout> | null = null;

	function dispatchChanged(page: number) {
		dispatch('changed', page);
	}

	function onNextButton() {
		let next = current + 1;
		if (next > total) {
			next = total;
		} else if (next < 1) {
			next = 1;
		}
		dispatchChanged(next);
	}

	function onPrevButton() {
		let prev = current - 1;
		if (prev > total) {
			prev = total;
		} else if (prev < 1) {
			prev = 1;
		}
		dispatchChanged(prev);
	}

	function onPageInput() {
		if (pageInputTimeoutID) {
			clearTimeout(pageInputTimeoutID);
		}

		pageInputTimeoutID = setTimeout(() => {
			const inputPageInt = parseInt(inputPage, 10);
			if (isNaN(inputPageInt)) {
				return;
			}

			const isBadInput = inputPageInt < 1 || inputPageInt === current || inputPageInt > total;
			if (isBadInput) {
				return;
			}

			dispatchChanged(inputPageInt);
		}, 1000);
	}
</script>

<div class="pagination">
	<div class="paginator">
		<div class="prev">
			{#if current > 1}
				<div class="prev__butt pointer center" on:click={onPrevButton} />
			{/if}
		</div>

		<div class="page">
			<input
				type="number"
				placeholder="page"
				bind:value={inputPage}
				on:input={onPageInput}
			/>
		</div>

		<div class="next">
			{#if current < total}
				<div class="next__butt pointer center" on:click={onNextButton} />
			{/if}
		</div>
	</div>

	<div class="total">
		<div class="count center">{total}</div>
	</div>
</div>

<style lang="scss">
	@import '../../lib/assets/vars.scss';

	input[type='number'] {
		-moz-appearance: textfield;
	}

	.center {
		display: flex;
		justify-content: center;
		align-items: center;
		height: 100%;
		width: 100%;
	}

	.pagination {
		width: 100%;
		height: var(--paginator-height);
		display: grid;
		grid-template-columns: 1fr;
		grid-template-rows: 50% 1fr;
		gap: 12px;
		.paginator,
		.total {
			width: 100%;
			background: var(--color-level-1);
			border-radius: var(--border-radius);
			border: var(--color-border) 1px solid;
		}
		.paginator {
			display: grid;
			grid-template-rows: 1fr;
			grid-template-columns: repeat(3, 1fr);
			.prev,
			.next {
				&__butt {
					background: var(--color-level-2);
					width: 100%;
					height: 100%;
				}
			}
			.prev {
				&__butt {
					border-radius: var(--border-radius) 0 0 var(--border-radius);
				}
			}
			.page {
				height: 100%;
				width: 100%;
				display: flex;
				align-items: center;
				justify-content: center;
				input {
					width: 75%;
					height: 75%;
					background: var(--color-level-2);
					border: none;
					font-size: 1.4rem;
					text-align: center;
					text-indent: 0;
				}
			}
			.next {
				&__butt {
					border-radius: 0 var(--border-radius) var(--border-radius) 0;
				}
			}
		}
		.total {
			font-size: 1.3rem;
			padding: 12px;
		}
	}
</style>
