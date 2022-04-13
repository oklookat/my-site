<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	// files
	import type { Data } from '$lib_elven/types';
	import type { File, Params } from '$lib_elven/types/files';
	import NetworkFile from '$lib_elven/network/network_file';
	import Pagination from '$lib_elven/components/pagination.svelte';
	import FilesList from '$lib_elven/components/files_list.svelte';
	import Store from '$lib_elven/tools/store';
	import { browser } from '$app/env';

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
	let items: Data<File>;

	let container: HTMLDivElement;

	/** store unsubs */
	let unsub = {
		onSelected: undefined
	};

	function initStore() {
		Store.file.withSelectOption.set(true);
		unsub.onSelected = Store.file.selected.subscribe((v) => {
			if (!v) {
				return;
			}
			onSelect(v);
		});
	}

	function destroyStore() {
		unsub.onSelected();
		Store.file.withSelectOption.set(false);
		Store.file.selected.set(null);
	}

	onMount(async () => {
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
		if (p.page < 1) {
			p.page = 1;
		}
		try {
			const result = await networkFile.getAll(p);
			items = result;
		} catch (err) {}
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		params.page = page;
		await getAll();
	}

	/** on select button clicked in actions */
	function onSelect(file: File) {
		dispatch('selected', file);
	}

	function onClosed() {
		dispatch('closed');
	}
</script>

<div class="overlay base__overlay" bind:this={container}>
	<div class="overlay__main">
		<div class="close pointer" on:click={onClosed}>
			<svg viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg"
				><path
					d="M289.94,256l95-95A24,24,0,0,0,351,127l-95,95-95-95A24,24,0,0,0,127,161l95,95-95,95A24,24,0,1,0,161,385l95-95,95,95A24,24,0,0,0,385,351Z"
				/></svg
			>
		</div>

		<div class="base__container">
			<div class="files">
				<FilesList {items} onDeleted={() => {}} />
			</div>

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

			.close {
				width: 100%;
				height: 48px;
				background-color: var(--color-level-1);
				display: flex;
				justify-content: center;
				align-items: center;
				svg {
					fill: red;
					width: 30px;
					height: 30px;
				}
			}
			.files {
				height: 100%;
			}
			.files,
			.pages {
				width: 96%;
				margin: auto;
			}
			.pages {
				padding-top: 14px;
			}
		}
	}
</style>
