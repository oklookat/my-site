<script context="module" lang="ts">
	export const load: Load = async (e) => {
		let requestParams = new Params<File>('file', e.url.searchParams);

		let resp: Response | null = null;
		let items: Items<File>;

		const stuff = e.stuff;
		stuff.title = t.get('elven.files.title');

		const networkFile = new NetworkFile(getTokenFromSession(e));

		try {
			resp = await networkFile.getAll(requestParams.toObject(), e.fetch);
			if (resp.ok) {
				items = (await resp.json()) as Items<File>;
			}
		} catch (err) {
			throw Error(resp?.statusText);
		}

		return {
			status: resp?.status || 200,
			stuff: stuff,
			props: {
				items: items,
				params: requestParams
			}
		};
	};
</script>

<script lang="ts">
	import type { Items } from '$lib/types';
	import Pagination from '$lib/components/pagination.svelte';
	import type { File } from '$lib/types/files';
	import FilesToolbars from '$lib/components/elven/files_toolbars.svelte';
	import FilesList from '$lib/components/elven/files_list.svelte';
	import NetworkFile from '$lib/elven/entities/file/network';
	import {
		HandleRouteParam,
		Params,
		Refresh,
		type RPH_Data,
		type RPH_Event
	} from '$lib/tools/params';
	import type { Load } from '@sveltejs/kit';
	import { t } from '$lib/locale';
	import { getTokenFromSession } from '$lib/tools';

	/** files data */
	export let items: Items<File>;

	/** request params */
	export let params: Params<File>;

	async function onPageChanged(page: number) {
		await onParamChanged({ name: 'page', val: page });
	}

	async function onUploaded() {
		await onParamChanged({ name: 'page', val: 1 });
	}

	/** on request param changed */
	async function onParamChanged(event: RPH_Event<File>) {
		const data: RPH_Data<File> = {
			params: params,
			items: items
		};
		await HandleRouteParam<File>(event, data);
	}

	async function refresh() {
		const getData = async () => {
			return Promise.resolve({
				items: items,
				params: params
			});
		};
		await Refresh<File>(getData);
	}
</script>

<div class="files base__container">
	<FilesToolbars
		bind:params
		on:uploaded={async () => await onUploaded()}
		on:paramChanged={async (e) => await onParamChanged(e.detail)}
	/>

	<FilesList {items} on:deleted={async () => await refresh()} />

	<div class="pages">
		{#if items && items.meta}
			<Pagination
				total={items.meta.total_pages}
				current={params.getParam('page')}
				on:changed={async (e) => await onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</div>
