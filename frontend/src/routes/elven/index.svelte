<script lang="ts">
	import { session } from '$app/stores';
	
</script>

<div class="hello base__container">
	<div class="hi">{$session.user.username}</div>
</div>

<style lang="scss">
	.hello {
		display: flex;
		flex-direction: column;
		flex-wrap: wrap;
		gap: 24px;
		height: 100%;
		border-radius: var(--border-radius);
		.hi {
			font-size: 2rem;
		}
	}
</style>
