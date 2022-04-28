<script lang="ts">
	export let percents: number;
</script>

<div class="progress" style="width: {percents}%" />

<style>
	.progress {
		border-radius: var(--border-radius);
		height: 100%;
		width: 0;
		background-color: var(--color);
	}
</style>
