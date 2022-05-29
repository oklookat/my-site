<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<{ click: MouseEvent }>();
</script>

<svg on:click={(e) => dispatch('click', e)} viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
	<rect id="view-box" fill="none" />
	<path
		d="M.22,10.22A.75.75,0,0,0,1.28,11.28l5-5a.75.75,0,0,0,0-1.061l-5-5A.75.75,0,0,0,.22,1.28l4.47,4.47Z"
		transform="translate(14.75 17.75) rotate(180)"
	/>
</svg>

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
