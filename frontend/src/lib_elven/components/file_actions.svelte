<script lang="ts">
	import { onMount } from 'svelte';
	import Overlay from '$elven/components/overlay.svelte';
	import Popup from '$elven/components/popup.svelte';
	import Store from '$elven/tools/store';
	import type { File } from '$elven/types/file';
	import NetworkFile from '$elven/network/file';
	import Preview from '$elven/components/preview.svelte';
	import { getUploadsURL, isTouchDevice } from '$elven/tools';
	import { t } from '$lib/locale';

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
		const isDelete = await window.$confirm($t('elven.files.deleteQuestion'));
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
			message = $t('elven.files.linkCopied');
		} catch (err) {
			message = $t('elven.files.notHavePermError');
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
		{#if withSelectOption}
			<div
				on:click={() => {
					Store.files.selected.set(file);
				}}
			>
				{$t('elven.files.select')}
			</div>
		{/if}

		<div on:click={() => onPreview(true)}>{$t('elven.files.preview')}</div>

		<div on:click={() => copyLink()}>{$t('elven.files.copyLink')}</div>

		{#if onDeleted}
			<div on:click={() => deleteFile()}>{$t('elven.files.delete')}</div>
		{/if}
	</svelte:component>
{/if}
