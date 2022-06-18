<script lang="ts">
	import Toolbar from '$lib/components/toolbar.svelte';
	import SearchBar from '$lib/components/search_bar.svelte';
	import { Start, type File } from '$elven/types/file';
	import FilesUploader from '$elven/components/files_uploader.svelte';
	import { createEventDispatcher, onDestroy } from 'svelte';
	import Store from '$elven/tools/store';
	import type { Params, RPH_Event } from '$elven/tools/params';
	

	/** request params */
	export let params: Params<File>;

	let searchValue = params.getParam('filename') || '';

	const dispatch = createEventDispatcher<{
		/** on request param changed */
		paramChanged: RPH_Event<File>;

		/** on file uploaded */
		uploaded: void;
	}>();

	const onUploadedExistsUnsub = Store.files.uploadedExists.subscribe((v) => {
		if (!v) {
			return;
		}
		searchValue = v.original_name;
	});

	onDestroy(() => {
		onUploadedExistsUnsub();
		Store.files.uploadedExists.set(null);
	});

	/** set 'start' param */
	function setStart(start: Start = Start.newest) {
		dispatch('paramChanged', { name: 'start', val: start });
	}

	/** search by filename */
	function search(val: string) {
		dispatch('paramChanged', { name: 'filename', val: val });
	}
</script>

<div class="toolbars">
	<FilesUploader on:uploaded={() => dispatch('uploaded')} />

	<div class="search">
		<SearchBar bind:value={searchValue} on:search={(e) => search(e.detail)} />
	</div>

	<Toolbar>
		{#if params}
			{#if params.getParam('start') === Start.newest}
				<div class="item" on:click={() => setStart(Start.oldest)}>newest</div>
			{:else if params.getParam('start') === Start.oldest}
				<div class="item" on:click={() => setStart(Start.newest)}>oldest</div>
			{/if}
		{/if}
	</Toolbar>
</div>

<style lang="scss">
	.toolbars {
		display: flex;
		flex-direction: column;
		gap: 12px;
		width: 100%;
		.search {
			height: 54px;
		}
	}
</style>
