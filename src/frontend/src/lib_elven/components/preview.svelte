<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import type { FileTypeSelector } from '$lib_elven/tools/extension';

	/** on preview closed */
	export let onClose: () => void;

	/** file url */
	export let url: URL;

	/** file extension */
	export let extension: FileTypeSelector;

	let isImage = false;
	let isVideo = false;
	let isAudio = false;
	let isSupported = false;

	onMount(() => {
		// check support
		isImage = extension.selected === 'IMAGE';
		isVideo = extension.selected === 'VIDEO';
		isAudio = extension.selected === 'AUDIO';
		isSupported = isImage || isVideo || isAudio;
		if (!isSupported) {
			window.$notify?.add({ message: 'Unsupported.' });
			onClose();
			return;
		}

		document.body.classList.add('no-scroll');
	});

	onDestroy(() => {
		document.body.classList.remove('no-scroll');
	});

	/** play audio by url */
	function playAudio() {
		if (!window.$player) {
			onClose();
			return;
		}
		window.$player.clearPlaylist();
		window.$player.addToPlaylist(url.toString());
		window.$player.play();
		onClose();
	}
</script>

<div class="preview base__overlay" on:click|self={onClose}>
	{#if isAudio}
		{playAudio()}
	{/if}
	<div class="watchable">
		{#if isImage}
			<img decoding="async" loading="lazy" src={url.toString()} alt="preview" />
		{:else if isVideo}
			<video controls src={url.toString()}>
				<track default kind="captions" srclang="en" src="" />
			</video>
		{/if}
	</div>
</div>

<style lang="scss">
	.preview {
		.watchable {
			padding: 12px;
			max-width: 90%;
			max-height: 90%;
			img,
			video {
				width: 100%;
				height: 100%;
			}
		}
	}
</style>
