<div class="toolbar base__toolbar">
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
			height: max-content;
			display: flex;
			flex-direction: row;
			gap: 12px;
			padding: 12px;
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
