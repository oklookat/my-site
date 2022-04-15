<script lang="ts">
	// tools
	// ui
	import Toolbar from '$lib_elven/components/toolbar.svelte';
	import SearchBar from '$lib_elven/components/search_bar.svelte';
	// files
	import type { Params } from '$lib_elven/types/files';
	import { Start } from '$lib_elven/types/files';
	import FilesUploader from '$lib_elven/components/files_uploader.svelte';
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	import Store from '$lib_elven/tools/store';

	/** request params */
	export let params: Params;

	/** refresh files */
	export let refresh: () => Promise<void>;

	const dispatch = createEventDispatcher<{
		/** on request param changed */
		paramChanged: { name: string; val: string };
	}>();

	let existsSearchValue = undefined;
	const onUploadedExistsUnsub = Store.onUploadedFileExists.subscribe((v) => {
		if (!v) {
			return;
		}
		existsSearchValue = v.original_name;
	});

	onDestroy(() => {
		onUploadedExistsUnsub();
	});

	onMount(() => {
		// if we have initial filename param, set it to searchbar
		if (params.filename) {
			existsSearchValue = params.filename;
		}
	});

	/** set 'start' param */
	function setStart(start: Start = Start.newest) {
		params.start = start;
		dispatch('paramChanged', { name: 'start', val: start });
	}

	/** search by filename */
	function search(val: string) {
		params.filename = val;
		dispatch('paramChanged', { name: 'filename', val: val });
	}

	/** on file uploaded */
	async function onUploaded() {
		params.page = 1;
		await refresh();
	}
</script>

<div class="toolbars">
	<FilesUploader on:uploaded={async () => await onUploaded()} />

	<div class="oneline">
		<div class="sort">
			<Toolbar>
				<div class="sort-by-old">
					{#if params && params.start === Start.newest}
						<div class="item" on:click={() => setStart(Start.oldest)}>newest</div>
					{/if}
					{#if params && params.start === Start.oldest}
						<div class="item" on:click={() => setStart(Start.newest)}>oldest</div>
					{/if}
				</div>
			</Toolbar>
		</div>
		<div class="search">
			<SearchBar
				value={existsSearchValue}
				on:search={(e) => search(e.detail)}
				placeholder="search"
			/>
		</div>
	</div>
</div>

<style lang="scss">
	.toolbars {
		display: flex;
		flex-direction: column;
		gap: 12px;
		width: 100%;
		.oneline {
			display: flex;
			gap: 14px;
			width: 100%;
			.sort {
				width: 50%;
			}
			.search {
				height: 54px;
				width: 50%;
			}
		}
	}
</style>
