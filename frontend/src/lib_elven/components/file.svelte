<script lang="ts">
	import { dateToReadable } from '$elven/tools/dates';
	import Extension from '$elven/tools/extension';
	import type { File } from '$elven/types/file';
	import FileActions from '$elven/components/file_actions.svelte';
	import { bytesToReadable, cutString, getUploadsWith } from '$elven/tools';

	export let file: File;
	$: convert(file);

	/** on file deleted */
	export let onDeleted: () => void;

	/** convert file path, extension etc */
	function convert(file: File) {
		if (!file) {
			return;
		}
		if (!(file.pathConverted instanceof URL)) {
			file.pathConverted = getUploadsWith(file.path);
		}
		if (!file.extensionsSelector) {
			file.extensionsSelector = Extension.getSelector(file.extension);
		}
		if (!file.sizeConverted) {
			file.sizeConverted = bytesToReadable(file.size);
		}
		if (!file.createdAtConverted) {
			file.createdAtConverted = dateToReadable(file.created_at);
		}
		if (!file.original_name_short) {
			file.original_name_short = cutString(file.original_name);
		}
	}

	/** is file selected? (actions menu/overlay opened) */
	let isSelected = false;
	let selectedMouseEvent: MouseEvent;
	function onSelected(e: MouseEvent) {
		selectedMouseEvent = e;
		isSelected = true;
	}
</script>

{#if isSelected}
	<FileActions
		{file}
		mouseEvent={selectedMouseEvent}
		onDisabled={() => (isSelected = false)}
		onDeleted={() => onDeleted()}
	/>
{/if}

<div class="file" on:click={(e) => onSelected(e)}>
	{#if file.extensionsSelector && file.extensionsSelector.selected === 'IMAGE' && file.pathConverted}
		<div class="preview">
			<img decoding="async" loading="lazy" src={file.pathConverted.href} alt="" />
		</div>
	{:else}
		<div class="preview unknown" />
	{/if}

	<div class="meta">
		<div class="title">{file.original_name_short}</div>
		<div class="info">
			<div class="created">{file.createdAtConverted}</div>
			<div class="size">{file.sizeConverted}</div>
		</div>
	</div>
</div>

<style lang="scss">
	.file {
		cursor: pointer;
		border-radius: var(--border-radius);
		background-color: var(--color-level-1);
		font-size: 1rem;
		padding: 12px;
		width: 100%;
		min-height: 54px;
		display: flex;
		flex-direction: row;
		gap: 12px;

		.preview {
			max-height: 54px;
			width: 50px;
			display: flex;
			justify-content: center;
			align-items: center;

			* {
				height: 100%;
				width: 50px;
				object-fit: cover;
			}

			&.unknown {
				background-color: var(--color-level-2);
				border-radius: var(--border-radius);
				padding: 8px;
			}
		}

		.meta {
			display: flex;
			flex-direction: column;
			gap: 12px;

			.title {
				width: 100%;
			}

			.info {
				color: var(--color-text-inactive);
				display: flex;
				flex-direction: row;
				gap: 14px;
				padding-top: auto;
			}
		}
	}
</style>
