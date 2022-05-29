<script lang="ts">
	import { onMount } from 'svelte';
	import type { FileTypeSelector } from '$elven/tools/extension';

	import Overlay from '$lib/components/overlay.svelte';

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
			window.$notify?.add({ message: 'Unsupported' });
			onClose();
			return;
		}

		if (isAudio) {
			playAudio();
			return;
		}
	});

	/** play audio by url */
	function playAudio() {
		window.$player?.clearPlaylist();
		window.$player?.addToPlaylist(url);
		window.$player?.play();
		onClose();
	}
</script>

<Overlay {onClose}>
	<div class="preview">
		<div class="watchable">
			{#if isImage}
				<img decoding="async" loading="lazy" src={url.toString()} alt="" />
			{:else if isVideo}
				<video controls src={url.toString()}>
					<track default kind="captions" srclang="en" src="" />
				</video>
			{/if}
		</div>
	</div>
</Overlay>

<style lang="scss">
	.preview {
		align-self: center;
		justify-self: center;
		width: 80%;
		max-width: 700px;
		.watchable {
			width: 100%;
			display: flex;
			justify-content: center;
			img,
			video {
				width: 100%;
			}
		}
	}
</style>
