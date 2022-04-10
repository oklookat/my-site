<script lang="ts">
	import CFile from '$lib_elven/components/file.svelte';

	import type { Data } from '$lib_elven/types';
	import type { File } from '$lib_elven/types/files';

	export let items: Data<File>;
	export let onSelected: (file: File, counter: number, e: MouseEvent) => void;

    function select(file: File, counter: number, e: MouseEvent) {
        onSelected(file, counter, e)
    }
</script>

<div class="list">
	{#if items && items.data}
		{#each Object.entries(items.data) as [counter, file]}
			<CFile {file} onSelected={(e) => select(file, parseInt(counter, 10), e)} />
		{/each}
	{/if}
</div>