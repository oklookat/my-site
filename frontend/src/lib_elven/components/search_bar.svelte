<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { t } from '$lib/locale';

	export let placeholder = $t('elven.components.searchbar.search');

	export let value = '';

	let isInitial = true;

	$: onValueChanged(value);
	function onValueChanged(val: string) {
		if (isInitial) {
			isInitial = false;
			return;
		}

		if (!val) {
			val = '';
		}

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
				if (!inputEL) {
					return;
				}
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
		bind:value
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
