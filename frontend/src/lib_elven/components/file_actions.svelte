<script lang="ts">
	import { onMount } from 'svelte';
	import ContextMenu from '$lib/components/context_menu.svelte';
	import Store from '$elven/tools/store';
	import type { File } from '$elven/types/file';
	import NetworkFile from '$elven/network/file';
	import Preview from '$elven/components/preview.svelte';
	import { getUploadsURL, isTouchDevice } from '$elven/tools';
	import OverlayMobile from '$lib/components/overlay_mobile.svelte';

	/** itself */
	export let file: File;

	/** click that triggered this component */
	export let mouseEvent: MouseEvent;

	/** on deleted */
	export let onDeleted: () => void;

	/** on actions closed */
	export let onDisabled: () => void;

	/** add 'select' option to actions? */
	let withSelectOption = false;
	Store.files.withSelectOption.subscribe((v) => (withSelectOption = v))();

	/** choose component to render: overlay / popup */
	const render: {
		active: boolean;
		type: 'context' | 'overlay';
		component: typeof ContextMenu | typeof OverlayMobile | undefined;
		props: any;
	} = { active: true, type: 'overlay', component: undefined, props: null };

	/** is device with touchscreen? */
	let isTouch = false;

	onMount(() => {
		isTouch = isTouchDevice();
		if (isTouch) {
			render.type = 'overlay';
			render.component = OverlayMobile;
			render.props = {
				title: file.original_name_short,
				onClose: () => onDisabled()
			};
			return;
		}
		render.type = 'context';
		render.component = ContextMenu;
		render.props = {
			mouseEvent: mouseEvent,
			onDisabled: () => onDisabled()
		};
	});

	/** delete file */
	async function deleteFile() {
		if (!window.$confirm) {
			return;
		}
		onDisabled();
		const isDelete = await window.$confirm('Delete file');
		if (!isDelete) {
			return;
		}
		try {
			const resp = await NetworkFile.delete(file.id);
			if (resp.ok) {
				onDeleted();
			}
		} catch (err) {}
	}

	/** copy link to clipboard */
	async function copyLink() {
		let message = '';
		const path = file.path;
		const formattedPath = getUploadsURL().toString() + `/${path}`;
		try {
			await navigator.clipboard.writeText(formattedPath);
			message = 'Link copied';
		} catch (err) {
			message = 'Not have permission';
		}
		window.$notify?.add({ message });
		onDisabled();
	}

	let isPreviewActive = false;

	/** preview image/video/etc */
	function onPreview(enabled: boolean) {
		if (!enabled) {
			onDisabled();
			isPreviewActive = false;
			return;
		}
		const isValid = !!(file.pathConverted && file.extensionsSelector);
		if (!isValid) {
			return;
		}
		render.active = false;
		isPreviewActive = true;
	}
</script>

{#if isPreviewActive && file.pathConverted}
	<Preview
		onClose={() => onPreview(false)}
		url={file.pathConverted}
		extension={file.extensionsSelector}
	/>
{/if}

{#if render.active}
	<svelte:component this={render.component} {...render.props}>
		<div class="actions {render.type}">
			{#if withSelectOption}
				<div
					on:click={() => {
						Store.files.selected.set(file);
					}}
				>
					select
				</div>
			{/if}

			<div on:click={() => onPreview(true)}>preview</div>

			<div on:click={() => copyLink()}>copy link</div>

			<div on:click={() => deleteFile()}>delete</div>
		</div>
	</svelte:component>
{/if}

<style lang="scss">
	.actions {
		height: 100%;
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		> * {
			cursor: pointer;
			&:hover {
				background-color: var(--color-hover);
			}
			display: flex;
			align-items: center;
			justify-content: center;
		}
		&.overlay,
		&.context {
			> * {
				width: 100%;
			}
		}
		&.overlay {
			padding-top: 12px;
			gap: 14px;
			// item
			> * {
				height: 64px;
				background-color: var(--color-level-1);
			}
		}
		&.context {
			// item
			> * {
				height: 44px;
			}
		}
	}
</style>
