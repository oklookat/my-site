<script lang="ts">
	import NetworkFile from '$lib_elven/network/network_file';
	import Store from '$lib_elven/tools/store';
	import type { File as TFile } from '$lib_elven/types/files';

	import { createEventDispatcher } from 'svelte';
	// files
	const dispatch = createEventDispatcher<{
		/** on file uploaded */
		uploaded: void;
	}>();

	/** allow upload multiple files? */
	export let multipleUploading = false;

	/** file input for upload */
	let inputEL;

	/** text used in upload field */
	let hintText = '';

	/** drag on uploader? */
	let isDragActive = false;

	/** is file(s) uploads on server now? */
	let isUploadingNow = false;

	/** set hint text */
	dragSwitcher(false);

	/** when file upload button clicked */
	function onUploadClick() {
		if (!inputEL) {
			return;
		}
		inputEL.value = '';
		inputEL.click();
	}

	/** when file changed on file input */
	async function onInputChange(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files.length < 1) {
			return 0;
		}
		const file = target.files[0];
		upload(file);
	}

	/** uploader drag start / leave */
	function onDrag(e: DragEvent) {
		e.preventDefault();
		if (isUploadingNow) {
			return;
		}
		const isStart = e.type === 'dragenter' && !isDragActive;
		if (isStart) {
			dragSwitcher(true);
		} else if (e.type === 'dragleave' && isDragActive) {
			dragSwitcher(false);
		}
	}

	/** set isDragActive & hint text */
	function dragSwitcher(enable: boolean) {
		if (enable) {
			hintText = 'release mouse to upload';
		} else {
			hintText = 'click or drag to upload';
		}
		isDragActive = enable;
	}

	/** https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/File_drag_and_drop */
	function onDrop(e: DragEvent) {
		e.preventDefault();
		if (isUploadingNow) {
			return;
		}
		dragSwitcher(false);
		for (let i = 0; i < e.dataTransfer.items.length; i++) {
			const isFile = e.dataTransfer.items[i].kind === 'file';
			if (!isFile) {
				continue;
			}
			const file = e.dataTransfer.items[i].getAsFile();
			upload(file);
			// upload only one if multipleUploading disabled
			if (!multipleUploading) {
				break;
			}
		}
	}

	/** https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/File_drag_and_drop */
	function onDragOver(e: DragEvent) {
		e.preventDefault();
	}

	async function upload(file: File) {
		if (!file || isUploadingNow) {
			return;
		}
		isUploadingNow = true;
		try {
			// use empty token because we have token in cookie
			const resp = await NetworkFile.upload(file);
			if (resp.ok) {
				dispatch('uploaded');
			} else if(resp.status === 409) {
				// if uploaded file already exists
				const theFile = await resp.json() as TFile
				Store.onUploadedFileExists.set(theFile)
			}
			
		} catch (err) {
		} finally {
			isUploadingNow = false;
		}
	}
</script>

<div
	class="uploader base__toolbar--big"
	on:click={() => onUploadClick()}
	on:dragenter={onDrag}
	on:dragleave={onDrag}
	on:drop={onDrop}
	on:dragover={onDragOver}
>
	<div class="hint content">{hintText}</div>
	<input
		type="file"
		style="display: none"
		multiple={multipleUploading}
		bind:this={inputEL}
		on:input={onInputChange}
	/>
</div>

<style lang="scss">
	.uploader {
		cursor: pointer;
		width: 100%;
		height: 64px;
		display: flex;
		align-items: center;
		justify-content: center;
		.hint {
			pointer-events: none;
			opacity: 0.5;
		}
	}
</style>
