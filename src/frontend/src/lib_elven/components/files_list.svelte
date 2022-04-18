<script lang="ts">
	import CFile from '$lib_elven/components/file.svelte';

	import type { Items } from '$lib_elven/types';
	import type { File } from '$lib_elven/types/files';

	export let items: Items<File>;

	/** on file deleted */
	export let onDeleted: () => void;

	/** delete file from files array */
	function deleteFromArray(counter: number) {
		delete items.data[counter];
		items = items;
	}

	function whenDeleted(counter: number) {
		deleteFromArray(counter);
		onDeleted();
	}
</script>

<div class="list">
	{#if items && items.data}
		{#each Object.entries(items.data) as [counter, file]}
			<CFile {file} onDeleted={() => whenDeleted(parseInt(counter, 10))} />
		{/each}
	{/if}
</div>
