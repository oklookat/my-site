<script lang="ts">
	import { onMount } from 'svelte';

	export let path: string;

	onMount(() => {
		if (!path.startsWith('/')) {
			path = '/' + path;
		}
	});
</script>

<a href="/elven{path}"><slot /></a>
