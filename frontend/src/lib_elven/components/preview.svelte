<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import type { FileTypeSelector } from '$elven/tools/extension';
	import { t } from '$lib/locale';
	import { toggleBodyScroll } from '$elven/tools';

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

	let setDefScroll: () => void;
	onMount(() => {
		setDefScroll = toggleBodyScroll();

		// check support
		isImage = extension.selected === 'IMAGE';
		isVideo = extension.selected === 'VIDEO';
		isAudio = extension.selected === 'AUDIO';
		isSupported = isImage || isVideo || isAudio;
		if (!isSupported) {
			window.$notify?.add({ message: $t('elven.general.unsupported') });
			onClose();
			return;
		}

		if (isAudio) {
			playAudio();
			return;
		}
	});

	onDestroy(() => {
		setDefScroll();
	});

	/** play audio by url */
	function playAudio() {
		if (!window.$player) {
			onClose();
			return;
		}
		window.$player.clearPlaylist();
		window.$player.addToPlaylist(url);
		window.$player.play();
		onClose();
	}
</script>

<div class="preview base__overlay" on:click|self={onClose}>
	<div class="watchable">
		{#if isImage}
			<img decoding="async" loading="lazy" src={url.toString()} alt={$t('elven.files')} />
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
			width: 80%;
			max-width: 700px;
			display: flex;
			justify-content: center;

			img,
			video {
				width: 100%;
				max-width: 844px;
			}
		}
	}
</style>
