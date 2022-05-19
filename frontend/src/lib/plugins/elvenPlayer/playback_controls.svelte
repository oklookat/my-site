<script lang="ts">
	import Next from './icons/next.svelte';
	import Prev from './icons/prev.svelte';
	import Play from './icons/play.svelte';
	import { audioComponent, isPlaying } from './store';
	import Pause from './icons/pause.svelte';

	function play() {
		window.$player?.play();
	}

	function pause() {
		window.$player?.pause();
	}

	function next() {
		window.$player?.next();
	}

	function prev() {
		window.$player?.prev();
	}
</script>

<div class="controls">
	<div class="prev" on:click={prev}>
		<Prev />
	</div>

	{#if $isPlaying && $audioComponent}
		<div class="pause" on:click={pause}>
			<Pause />
		</div>
	{:else if !$isPlaying || !$audioComponent}
		<div class="play" on:click={play}>
			<Play />
		</div>
	{/if}

	<div class="next" on:click={next}>
		<Next />
	</div>
</div>

<style lang="scss">
	.controls {
		width: max-content;
		height: max-content;
		display: grid;
		grid-template-rows: 1fr;
		grid-template-columns: repeat(3, auto);
		gap: 32px;
		align-items: center;
		justify-content: center;
		div {
			height: 100%;
			width: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
			cursor: pointer;
		}
	}
</style>
