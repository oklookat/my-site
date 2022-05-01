<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<{ click: MouseEvent }>();
</script>

<svg on:click={(e) => dispatch('click', e)} viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg"
	><path d="M41,14H7a2,2,0,0,1,0-4H41A2,2,0,0,1,41,14Z" /><path
		d="M41,26H7a2,2,0,0,1,0-4H41A2,2,0,0,1,41,26Z"
	/><path d="M41,38H7a2,2,0,0,1,0-4H41A2,2,0,0,1,41,38Z" /></svg
>

<style lang="scss">
	svg {
		fill: white;
		cursor: pointer;
		width: 100%;
		height: 100%;
	}
	@media screen and(prefers-color-scheme: light) {
		svg {
			fill: black;
		}
	}
</style>
