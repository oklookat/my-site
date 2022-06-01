<script lang="ts">
	import { browser } from '$app/env';

	import type { Settings } from '$elven/plugins/elvenProgress/types';

	import { onDestroy, onMount } from 'svelte';

	let percents: number;
	let settings: Settings = {
		height: '2px',
		basicLoading: {
			startTo: 45,
			startSpeed: 30,
			finishSpeed: 1
		}
	};

	onMount(() => {
		window.$progress = {
			get percents(): number {
				return percents;
			},
			set percents(val: number) {
				if (val > 100) {
					val = 100;
				} else if (val < 0) {
					val = 0;
				}
				percents = val;
			},
			startBasic() {
				const intervalID = setInterval(() => {
					if (percents < settings.basicLoading.startTo) {
						percents++;
						return;
					}
					clearInterval(intervalID);
				}, settings.basicLoading.startSpeed);
			},
			finishBasic() {
				percents = settings.basicLoading.startTo;
				const intervalID = setInterval(() => {
					if (percents < 100) {
						percents++;
						return;
					}
					this.reset();
					clearInterval(intervalID);
				}, settings.basicLoading.finishSpeed);
			},
			reset() {
				percents = 0;
			}
		};
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		window.$progress = undefined;
	});
</script>

<div class="progress">
	<div class="line" style={`width: ${percents}%`} />
</div>

<style lang="scss">
	.progress {
		z-index: 9999;
		cursor: default;
		top: 0;
		position: absolute;
		width: 100%;
		height: 2px;
		.line {
			height: 100%;
			width: 0;
			background-color: #a097dc;
		}
	}
</style>
