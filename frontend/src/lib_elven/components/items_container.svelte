<div class="items">
	<slot name="up" />
	<div class="list">
		<slot name="list" />
	</div>
	<div class="pages">
		<slot name="pages" />
	</div>
</div>

<style lang="scss">
	@import '../../lib/assets/vars.scss';

	.items {
		margin: auto;
		max-width: 744px;
		height: 100%;
		display: grid;
		grid-template-rows: max-content 1fr max-content;
		grid-template-columns: 1fr;
		gap: 14px;

		:global(.list > div) {
			//padding-top: 12px;
			width: 100%;
			height: 100%;
			min-height: 42px;

			display: grid;
			grid-template-columns: 1fr;
			grid-auto-rows: max-content;
			gap: 14px;

			@media screen and(max-width: $readable-max-width) {
				justify-content: center;
			}
		}

		:global(.pages > div) {
			width: 100%;
			height: max-content;
			display: flex;
			justify-content: center;

			@include media__desktop {
				display: block;
			}

			margin: auto;
		}
	}
</style>
