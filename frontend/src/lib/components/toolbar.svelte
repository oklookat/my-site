<div class="toolbar">
	<div class="content">
		<slot />
	</div>
</div>

<style lang="scss">
	.toolbar {
		height: max-content;
		display: flex;
		flex-direction: column;
		gap: 12px;
		.content {
			border: 1px solid var(--color-border);
			background-color: var(--color-level-1);
			min-height: 48px;
			height: max-content;
			display: flex;
			flex-wrap: wrap;
			align-items: center;
			flex-direction: row;
			gap: 12px;
			padding: 10px;
			font-size: 1rem;
			:global(> *) {
				cursor: pointer;
				height: 100%;
				display: flex;
				align-items: center;
				justify-content: center;
			}
		}
	}
</style>
