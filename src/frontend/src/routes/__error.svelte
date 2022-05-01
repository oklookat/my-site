<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';

	export const load: Load = async (event) => {
		return {
			props: {
				title: `${event.status}`
			}
		};
	};
</script>

<script lang="ts">
	export let title: string;
</script>

{title}
