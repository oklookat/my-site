<script lang="ts" context="module">
import type { Load } from "@sveltejs/kit";



export const load: Load = async (e) => {
	return {
		status: 302,
		redirect: "/elven/articles"
	}
}

</script>
