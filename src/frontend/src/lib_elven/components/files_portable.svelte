<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	// ui
	import Animation from '$lib_elven/tools/animation';
	// files
	import type { Data } from '$lib_elven/types';
	import type { File, Params } from '$lib_elven/types/files';
	import NetworkFile from '$lib_elven/network/network_file';
	import Pagination from '$lib_elven/components/pagination.svelte';
	import FileActions from '$lib_elven/components/file_actions.svelte';
	import FilesList from '$lib_elven/components/files_list.svelte';

	const dispatch = createEventDispatcher<{
		/** on 'select' option clicked on file */
		selected: File;
		/** on files closed */
		closed: void;
	}>();

	const networkFile = new NetworkFile('');

	/** request params from portable mode */
	export let params: Params = undefined;

	/** files loaded? */
	let loaded = false;

	/** response data */
	let items: Data<File>;

	let container: HTMLDivElement;

	/** is file selected? */
	let isSelected = false;

	/** selected file */
	let selected: {
		file: File;
		mouseEvent: MouseEvent;
	} = { file: null, mouseEvent: null };

	onMount(async () => {
		document.body.classList.add('no-scroll');
		Animation.fadeIn(container, 10);
		await getAll();
	});

	onDestroy(() => {
		document.body.classList.remove('no-scroll');
	});

	/** get all files */
	async function getAll(p: Params = params) {
		params = params;
		if (p.page < 1) {
			p.page = 1;
		}
		loaded = false;
		try {
			const result = await networkFile.getAll(p);
			items = result;
			loaded = true;
		} catch (err) {}
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		params.page = page;
		await getAll();
	}

	/** select file */
	function select(file: File, event: MouseEvent) {
		selected.file = file;
		selected.mouseEvent = event;
		isSelected = true;
	}

	/** on select button clicked in actions */
	function onSelect(file: File) {
		dispatch('selected', file);
	}

	function onClosed() {
		isSelected = false;
		dispatch('closed');
	}
</script>

{#if isSelected}
	<FileActions
		file={selected.file}
		mouseEvent={selected.mouseEvent}
		withSelect={true}
		onDisabled={() => (isSelected = false)}
		onSelectClicked={() => onSelect(selected.file)}
	/>
{/if}

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
				<FilesList {items} onSelected={(file, counter, e) => select(file, e)} />
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
