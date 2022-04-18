<script lang="ts">
	import { browser } from '$app/env';
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	// utils
	import Store from '$lib_elven/tools/store';
	// ui
	import Pagination from '$lib_elven/components/pagination.svelte';
	// files
	import type { Items } from '$lib_elven/types';
	import type { File, Params } from '$lib_elven/types/files';
	import NetworkFile from '$lib_elven/network/network_file';
	import FilesList from '$lib_elven/components/files_list.svelte';
	import FilesToolbars from '$lib_elven/components/files_toolbars.svelte';
	import ToolsFiles from '$lib_elven/tools/files';
	import Utils from '$lib_elven/tools';

	const dispatch = createEventDispatcher<{
		/** on 'select' option clicked on file */
		selected: File;

		/** on files closed */
		closed: void;
	}>();

	const networkFile = new NetworkFile('');

	/** request params from portable mode */
	export let params: Params = undefined;

	/** response data */
	let items: Items<File>;

	/** store unsubs */
	let unsub = {
		onSelected: undefined
	};

	function initStore() {
		Store.file.withSelectOption.set(true);
		unsub.onSelected = Store.file.selected.subscribe((file) => {
			if (!file) {
				return;
			}
			dispatch('selected', file);
		});
	}

	function destroyStore() {
		unsub.onSelected();
		Store.file.withSelectOption.set(false);
		Store.file.selected.set(null);
	}

	onMount(async () => {
		const defaultParams = ToolsFiles.getDefaultParams();
		if (!params) {
			params = defaultParams;
		} else {
			params = Object.assign(defaultParams, params);
		}
		initStore();
		document.body.classList.add('no-scroll');
		await getAll();
	});

	onDestroy(() => {
		destroyStore();
		if (!browser) {
			return;
		}
		document.body.classList.remove('no-scroll');
	});

	/** get all files */
	async function getAll(p: Params = params) {
		params = params;
		const backupPage = p.page;
		if (p.page < 1) {
			p.page = 1;
		}
		let isError = false;
		try {
			const resp = await networkFile.getAll(p);
			if (resp.status === 200) {
				items = await resp.json();
				return;
			}
			isError = true;
		} catch (err) {
			isError = true;
		}
		if (isError) {
			// revert page change
			p.page = backupPage;
		}
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		params.page = page;
		await getAll();
	}

	async function refresh() {
		const getPage = async () => {
			return params.page;
		};
		const setPage = async (newPage: number) => {
			params.page = newPage;
		};
		const fetchItems = async (initial: boolean) => {
			if(initial) {
				return items
			}
			await onPageChanged(await getPage());
			return items;
		};
		await Utils.refresh(getPage, setPage, fetchItems);
	}

	async function onParamChanged(event: { name: string; val: string }) {
		params.page = 1;
		params[event.name] = event.val;
		await getAll();
	}
</script>

<div class="overlay base__overlay">
	<div class="overlay__main">
		<div
			class="back"
			on:click={() => {
				dispatch('closed');
			}}
		>
			<svg width="24px" height="24px" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
				<rect id="view-box" width="24" height="24" fill="none" />
				<path
					d="M.22,10.22A.75.75,0,0,0,1.28,11.28l5-5a.75.75,0,0,0,0-1.061l-5-5A.75.75,0,0,0,.22,1.28l4.47,4.47Z"
					transform="translate(14.75 17.75) rotate(180)"
				/>
			</svg>

			<slot name="back-title" />
		</div>

		<div class="base__container">
			<FilesToolbars
				bind:params
				{refresh}
				on:paramChanged={async (e) => onParamChanged(e.detail)}
			/>

			<FilesList {items} onDeleted={async () => await refresh()} />

			<div class="pages">
				{#if items && items.meta}
					<Pagination
						total={items.meta.total_pages}
						current={items.meta.current_page}
						on:changed={(e) => onPageChanged(e.detail)}
					/>
				{/if}
			</div>
		</div>
	</div>
</div>

<style lang="scss">
	.overlay {
		z-index: 9990;
		background-color: var(--color-body);
		display: block;
		&__main {
			z-index: 9990;
			overflow: auto;
			width: 100%;
			height: 100%;

			display: grid;
			grid-template-rows: max-content 1fr auto;
			gap: 12px;

			.back {
				cursor: pointer;
				width: 100%;
				height: 48px;
				background-color: var(--color-level-1);
				display: flex;
				align-items: center;
				svg {
					width: 30px;
					height: 30px;
					fill: var(--color-text);
				}
			}
		}
	}
</style>
