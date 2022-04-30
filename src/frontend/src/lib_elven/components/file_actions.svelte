<script lang="ts">
	import { onMount } from 'svelte';
	import Overlay from '$lib_elven/components/overlay.svelte';
	import Popup from '$lib_elven/components/popup.svelte';
	import Store from '$lib_elven/tools/store';
	import type { File } from '$lib_elven/types/files';
	import NetworkFile from '$lib_elven/network/network_file';
	import Preview from '$lib_elven/components/preview.svelte';
	import { getUploadsURL, isTouchDevice } from '$lib_elven/tools';
	import { _ } from 'svelte-i18n';

	/** file itself */
	export let file: File;

	/** click on file mouse event */
	export let mouseEvent: MouseEvent;

	/** on file deleted */
	export let onDeleted: () => void;

	/** on actions closed */
	export let onDisabled: () => void;

	/** add 'select' option to actions? */
	let withSelectOption = false;
	Store.files.withSelectOption.subscribe((v) => (withSelectOption = v))();

	/** choose component to render: overlay / popup */
	let render: {
		isOverlay: boolean;
		component: any;
		props: any;
	} = { isOverlay: true, component: null, props: null };

	/** is device with touchscreen? */
	let isTouch = false;

	onMount(() => {
		isTouch = isTouchDevice();
		if (isTouch) {
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
		if (!window.$confirm || !onDeleted || !onDisabled) {
			return;
		}
		onDisabled();
		const isDelete = await window.$confirm($_('elven.components.fileActions.deleteQuestion'));
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
		const formattedPath = getUploadsURL().toString() + `/${path}`;
		try {
			await navigator.clipboard.writeText(formattedPath);
			message = $_('elven.general.linkCopied');
		} catch (err) {
			message = $_('elven.general.notHavePermError');
		}
		window.$notify?.add({ message });
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
		onDisabled();
		isPreview = false;
	}
</script>

{#if isPreview && file.pathConverted && file.extensionsSelector}
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
						Store.files.selected.set(file);
					}}
				>
					{$_('elven.general.select')}
				</div>
			{/if}

			<div on:click={() => onPreview(true)}>{$_('elven.general.preview')}</div>

			<div on:click={() => copyLink()}>{$_('elven.general.copyLink')}</div>

			{#if onDeleted}
				<div on:click={() => deleteFile()}>{$_('elven.general.delete')}</div>
			{/if}
		</div>
	</svelte:component>
{/if}
