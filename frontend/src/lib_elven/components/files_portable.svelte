<script lang="ts">
	import { browser } from '$app/env';
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	import Store from '$elven/tools/store';
	import Pagination from '$lib/components/pagination.svelte';
	import type { Items } from '$elven/types';
	import type { File } from '$elven/types/file';
	import NetworkFile from '$elven/network/file';
	import FilesList from '$elven/components/files_list.svelte';
	import FilesToolbars from '$elven/components/files_toolbars.svelte';
	import { Params, Refresh, type RPH_Event } from '$elven/tools/params';
	import type { Unsubscriber } from 'svelte/store';
	import { toggleBodyScroll } from '$elven/tools';

	const dispatch = createEventDispatcher<{
		/** on 'select' option clicked on file */
		selected: File;

		/** on files closed */
		closed: void;
	}>();

	const networkFile = new NetworkFile('');

	/** request params from portable mode */
	export let params: Params<File>;

	/** response data */
	let items: Items<File>;

	/** store unsubs */
	let unsubSelected: Unsubscriber;

	function initStore() {
		Store.files.withSelectOption.set(true);
		unsubSelected = Store.files.selected.subscribe((file) => {
			if (!file) {
				return;
			}
			dispatch('selected', file);
		});
	}

	function destroyStore() {
		if (unsubSelected) {
			unsubSelected();
		}
		Store.files.withSelectOption.set(false);
		Store.files.selected.set(null);
	}

	let setDefScroll: () => void;
	onMount(async () => {
		if (!browser) {
			return;
		}
		if (!params) {
			params = new Params<File>('file');
		}
		initStore();
		setDefScroll = toggleBodyScroll();
		await getAll();
	});

	onDestroy(() => {
		destroyStore();
		if (!browser) {
			return;
		}
		setDefScroll();
	});

	/** get all files */
	async function getAll() {
		params = params;
		const backupPage = params.getParam('page');
		if (backupPage < 1) {
			params.setParam('page', 1);
		}

		let isError = false;
		try {
			const resp = await networkFile.getAll(params.toObject());
			if (!resp.ok) {
				isError = true;
				return;
			}
			items = await resp.json();
		} catch (err) {
			isError = true;
		}

		if (isError) {
			// revert page change
			params.setParam('page', backupPage);
		}
	}

	async function onUploaded() {
		await onParamChanged({ name: 'page', val: 1 });
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		params.setParam('page', page);
		await getAll();
	}

	async function refresh() {
		const getData = async () => {
			await getAll();
			return Promise.resolve({
				items: items,
				params: params
			});
		};
		const data = await Refresh<File>(getData, false);
		items = data.items;
		params = data.params;
	}

	async function onParamChanged(event: RPH_Event<File>) {
		params.setParam('page', 1);
		params.setParam(event.name, event.val);
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

			<div class="back__title">
				<slot name="back-title" />
			</div>
		</div>

		<div class="base__container">
			<FilesToolbars
				bind:params
				on:paramChanged={async (e) => await onParamChanged(e.detail)}
				on:uploaded={async () => await onUploaded()}
			/>

			<FilesList {items} on:deleted={async () => await refresh()} />

			<div class="pages">
				{#if items && items.meta}
					<Pagination
						bind:total={items.meta.total_pages}
						bind:current={items.meta.current_page}
						on:changed={async (e) => await onPageChanged(e.detail)}
					/>
				{/if}
			</div>
		</div>
	</div>
</div>

<style lang="scss">
	.overlay {
		z-index: 7000;
		background-color: var(--color-body);
		display: block;
		&__main {
			z-index: 8000;
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
					fill: var(--color-text);
				}
			}
		}
	}
</style>
