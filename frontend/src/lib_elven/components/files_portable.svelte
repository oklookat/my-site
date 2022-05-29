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
	import Overlay from '$lib/components/overlay.svelte';
	import ItemsContainer from '$elven/components/items_container.svelte';
	import Back from '$lib/icons/back.svelte';

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

	onMount(async () => {
		if (!browser) {
			return;
		}
		if (!params) {
			params = new Params<File>('file');
		}
		initStore();
		await getAll();
	});

	onDestroy(() => {
		destroyStore();
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

	function onClose(e: MouseEvent) {
		dispatch('closed');
	}
</script>

<Overlay {onClose}>
	<div class="portable">
		<div class="back" on:click={onClose}>
			<div class="icon">
				<Back />
			</div>

			<slot name="back-title" />
		</div>

		<div class="content">
			<ItemsContainer>
				<div slot="up">
					<FilesToolbars
						bind:params
						on:uploaded={async () => await onUploaded()}
						on:paramChanged={async (e) => await onParamChanged(e.detail)}
					/>
				</div>

				<div slot="list">
					<FilesList {items} on:deleted={async () => await refresh()} />
				</div>

				<div slot="pages">
					{#if items && items.meta}
						<Pagination
							total={items.meta.total_pages}
							current={params.getParam('page')}
							on:changed={async (e) => await onPageChanged(e.detail)}
						/>
					{/if}
				</div>
			</ItemsContainer>
		</div>
	</div>
</Overlay>

<style lang="scss">
	.portable {
		overflow: auto;
		background-color: var(--color-body);
		width: 100%;
		height: 100%;

		display: grid;
		grid-template-rows: max-content 1fr auto;

		.back {
			cursor: pointer;
			width: 100%;
			height: 48px;
			background-color: var(--color-level-1);
			display: flex;
			align-items: center;
			.icon {
				height: 100%;
				width: 30px;
			}
		}

		.content {
			width: 100%;
			height: 100%;
			padding: 12px;
		}
	}
</style>
