<script lang="ts">
	import { onMount } from 'svelte';
	// ui
	import Overlay from '$lib_elven/components/overlay.svelte';
	import Popup from '$lib_elven/components/popup.svelte';
	// utils
	import { Env } from '$lib_elven/tools/paths';
	import Validator from '$lib_elven/validators';
	// file
	import type { File } from '$lib_elven/types/files';
	import NetworkFile from '$lib_elven/network/network_file';

	/** add 'select' option to actions */
	export let withSelect = false;

	/** file itself */
	export let file: File;

	/** click on file mouse event */
	export let mouseEvent: MouseEvent;

	/** on file deleted */
	export let onDeleted: () => void = undefined;

	/** on actions closed */
	export let onDisabled: () => void;

	/** on 'select' option clicked */
	export let onSelectClicked: () => void = undefined;

	/** is device with touchscreen? */
	let isTouchDevice = false;

	/** choose component to render: overlay / popup */
	let render: {
		isOverlay: boolean;
		component: any;
		props: any;
	} = { isOverlay: true, component: null, props: null };

	onMount(() => {
		isTouchDevice = Validator.isTouchDevice();
		if (isTouchDevice) {
			render.component = Overlay;
			render.props = {
				onClose: () => onDisabled()
			};
			return;
		}
		render.isOverlay = false;
		render.component = Popup;
		render.props = {
			mouseEvent: mouseEvent,
			onDisabled: () => onDisabled()
		};
	});

	/** delete file */
	async function deleteFile() {
		if (!onDeleted || !onDisabled) {
			return;
		}
		onDisabled();
		const isDelete = await window.$choose.confirm('delete file');
		if (!isDelete) {
			return;
		}
		try {
			await NetworkFile.delete(file.id);
			onDeleted();
		} catch (err) {}
	}

	/** play audio by url */
	function playAudio() {
		const url = file.pathConverted.href;
		window.$player.playlist = { position: 0, sources: [url] };
		window.$player.play();
	}

	/** copy link to clipboard */
	async function copyLink() {
		let message = '';
		const path = file.path;
		const formattedPath = Env.getUploads() + '/' + path;
		try {
			await navigator.clipboard.writeText(formattedPath);
			message = 'Link copied to clipboard.';
		} catch (err) {
			message = 'Error: clipboard access denied.';
		}
		window.$notify.add({ message });
		onDisabled();
	}
</script>

<svelte:component this={render.component} {...render.props}>
	<div class="base__items {render.isOverlay ? 'extended' : ''}">
		{#if withSelect}
			<div
				on:click={() => {
					onSelectClicked();
				}}
			>
				select
			</div>
		{/if}

		{#if file.extensionsSelector.selected === 'AUDIO'}
			<div on:click={() => playAudio()}>play</div>
		{/if}

		<div on:click={() => copyLink()}>copy link</div>

		{#if onDeleted}
			<div on:click={() => deleteFile()}>delete</div>
		{/if}
	</div>
</svelte:component>
