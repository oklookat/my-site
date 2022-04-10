<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let placeholder = 'search';

	export let value = undefined;
	$: onValueChanged(value);
	function onValueChanged(v) {
		if (!inputEL) {
			return;
		}
		if (!value) {
			v = '';
		}
		inputEL.value = v;
		save();
	}

	let inputEL: HTMLInputElement;

	const dipatch = createEventDispatcher<{ search: string }>();
	const save = saveThrottle();

	function saveThrottle() {
		let timer: NodeJS.Timeout;
		return () => {
			if (timer) {
				clearTimeout(timer);
			}
			timer = setTimeout(() => {
				dipatch('search', inputEL.value);
			}, 1000);
		};
	}
</script>

<div class="search">
	<input
		class="search__input"
		type="text"
		{placeholder}
		bind:this={inputEL}
		on:input={() => save()}
	/>
</div>

<style lang="scss">
	.search {
		min-width: 64px;
		height: 100%;
		&__input {
			width: 100%;
			height: 100%;
		}
	}
</style>
