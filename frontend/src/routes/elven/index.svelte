<script lang="ts">
	import { session } from '$app/stores';
	
</script>

<div class="hello base__container">
	<div class="hi">{$session.user.username}</div>
	<div class="things">
		<a class="card" href="/">main site</a>
		<a class="card" href="/elven/settings">settings</a>
	</div>
</div>

<style lang="scss">
	.card {
		height: 48px;
		width: max-content;
		background-color: var(--color-level-1);
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--border-radius);
		padding: 8px;
	}
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
		.things {
			display: flex;
			flex-direction: row;
			flex-wrap: wrap;
			gap: 14px;
		}
	}
</style>
