<script lang="ts">
	import Toolbar from '$lib_elven/components/toolbar.svelte';
	import SearchBar from '$lib_elven/components/search_bar.svelte';
	import type { Params } from '$lib_elven/types/files';
	import { Start } from '$lib_elven/types/files';
	import FilesUploader from '$lib_elven/components/files_uploader.svelte';
	import { createEventDispatcher, onDestroy } from 'svelte';
	import Store from '$lib_elven/tools/store';

	/** request params */
	export let params: Params;

	let searchValue = '';

	if (params.filename) {
		searchValue = params.filename;
	}

	const dispatch = createEventDispatcher<{
		/** on request param changed */
		paramChanged: { name: string; val: string };

		/** on file uploaded */
		uploaded: void;
	}>();

	const onUploadedExistsUnsub = Store.onUploadedFileExists.subscribe((v) => {
		if (!v) {
			return;
		}
		searchValue = v.original_name;
	});

	onDestroy(() => {
		onUploadedExistsUnsub();
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
		<SearchBar bind:value={searchValue} on:search={(e) => search(e.detail)} placeholder="search" />
	</div>

	<Toolbar>
		{#if params}
			{#if params.start === Start.newest}
				<div class="item" on:click={() => setStart(Start.oldest)}>newest</div>
			{:else if params.start === Start.oldest}
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
