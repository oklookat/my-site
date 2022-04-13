<script lang="ts">
	import { onMount } from 'svelte';
	// ui
	import Overlay from '$lib_elven/components/overlay.svelte';
	import Popup from '$lib_elven/components/popup.svelte';
	// utils
	import { Env } from '$lib_elven/tools/paths';
	import Validator from '$lib_elven/validators';
	import Store from '$lib_elven/tools/store';
	// file
	import type { File } from '$lib_elven/types/files';
	import NetworkFile from '$lib_elven/network/network_file';
	import Preview from '$lib_elven/components/preview.svelte';

	/** file itself */
	export let file: File;

	/** click on file mouse event */
	export let mouseEvent: MouseEvent;

	/** on file deleted */
	export let onDeleted: () => void = undefined;

	/** on actions closed */
	export let onDisabled: () => void;

	/** add 'select' option to actions? */
	let withSelectOption = false;
	Store.file.withSelectOption.subscribe((v) => (withSelectOption = v))();

	/** choose component to render: overlay / popup */
	let render: {
		isOverlay: boolean;
		component: any;
		props: any;
	} = { isOverlay: true, component: null, props: null };

	/** is device with touchscreen? */
	let isTouchDevice = false;

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

	/** is preview active */
	let isPreview = false;

	/** preview image/video/etc */
	function onPreview(enabled: boolean) {
		if (enabled) {
			isPreview = true;
			return;
		}
		isPreview = false;
		onDisabled();
	}
</script>

{#if isPreview}
	<Preview
		onClose={() => onPreview(false)}
		url={file.pathConverted}
		extension={file.extensionsSelector}
	/>
{/if}

{#if !isPreview}
	<svelte:component this={render.component} {...render.props}>
		<div class="base__items {render.isOverlay ? 'extended' : ''}">
			{#if withSelectOption}
				<div
					on:click={() => {
						Store.file.selected.set(file);
					}}
				>
					select
				</div>
			{/if}

			<div on:click={() => onPreview(true)}>preview</div>

			<div on:click={() => copyLink()}>copy link</div>

			{#if onDeleted}
				<div on:click={() => deleteFile()}>delete</div>
			{/if}
		</div>
	</svelte:component>
{/if}
