<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<{
		/** when item selected */
		selected: string;
	}>();

	/** default selected value */
	export let selected: string | number;

	/** selectable elements (options). Key = option value, value = option text */
	export let selectable: Record<string, string>;

	function onChange(e: Event) {
		const target = e.target as HTMLSelectElement;
		dispatch('selected', target.value);
	}
</script>

<select class="select" on:change={onChange} value={selected}>
	{#each Object.entries(selectable) as [value, text]}
		<option {value}>{text}</option>
	{/each}
</select>

<style lang="scss">
	.select {
		color: var(--color-text);
		background-color: var(--color-level-2);
		box-sizing: border-box;
		border-radius: 0.4rem;
		min-width: 94px;
		max-width: fit-content;
		height: 100%;
		padding: 4px;
	}
</style>
