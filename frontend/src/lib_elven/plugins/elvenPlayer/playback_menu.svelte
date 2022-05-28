<script lang="ts">
	import Slider from './slider/slider.svelte';
	import {
		currentTime,
		currentTimePercents,
		duration,
		durationPretty,
		setCurrentTime,
		bufferedPercents,
		volumePercents
	} from './store';
	import { convertPercentsToCurrentTime, getPretty } from './utils';
	import { onDestroy } from 'svelte';
	import PlaybackControls from './playback_controls.svelte';
	import Overlay from '$lib/components/overlay.svelte';

	export let onClose: () => void;

	let currentTimePretty = getPretty($currentTime);
	let currentTimeNew = 0;

	let isMouseDown = false;
	let isMouseDownInit = true;
	$: watchMouseDown(isMouseDown);
	function watchMouseDown(val: boolean) {
		if (isMouseDownInit) {
			isMouseDownInit = false;
			return;
		}
		if (!val) {
			$setCurrentTime = currentTimeNew;
		}
	}

	const unsub1 = currentTime.subscribe((v) => {
		if (isMouseDown) {
			return;
		}
		currentTimePretty = getPretty(v);
	});

	onDestroy(() => {
		unsub1();
	});

	let slidedPercents = 0;
	function onProgressSlide(percents: number) {
		slidedPercents = percents;
		currentTimeNew = convertPercentsToCurrentTime(slidedPercents, $duration);
		currentTimePretty = getPretty(currentTimeNew);
	}
</script>

<Overlay {onClose}>
	<div class="main">
		<div>
			<div class="status">
				<div class="progress">
					<div class="buffered" style="width: {$bufferedPercents}%;" />
					<div class="itself">
						<Slider
							percents={$currentTimePercents}
							on:slide={(e) => onProgressSlide(e.detail)}
							on:mouseDown={() => (isMouseDown = true)}
							on:mouseUp={() => (isMouseDown = false)}
						/>
					</div>
				</div>

				<div class="time">
					<div class="current">{currentTimePretty}</div>
					<div class="total">{$durationPretty}</div>
				</div>
			</div>

			<div class="control">
				<div class="volume">
					<Slider percents={$volumePercents} on:slide={(e) => ($volumePercents = e.detail)} />
				</div>
				<div class="playback">
					<PlaybackControls />
				</div>
			</div>
		</div>
	</div>
</Overlay>

<style lang="scss">
	.main {
		align-self: center;
		justify-self: center;
		user-select: none;
		background-color: var(--color-level-1);
		border-radius: var(--border-radius);
		height: 254px;
		width: 304px;
		@media screen and(max-width: 644px) {
			max-width: 304px;
			width: 95%;
		}
		padding: 24px;

		> div {
			width: 100%;
			height: 100%;
			display: flex;
			flex-direction: column;
			gap: 15%;

			.status {
				display: flex;
				flex-direction: column;
				gap: 14px;

				.progress {
					height: 16px;
					position: relative;
					display: flex;
					flex-direction: row;
					.itself {
						height: 100%;
						width: 100%;
						position: absolute;
						z-index: 4;
					}
					.buffered {
						background-color: var(--color-level-2);
						height: 100%;
					}
				}

				.time {
					display: flex;
					flex-direction: row;
					justify-content: center;
					height: max-content;
					width: 100%;
					.current {
						flex: 1;
					}
				}
			}

			.control {
				align-self: center;
				height: max-content;
				width: max-content;
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				gap: 64px;
				.volume {
					background-color: var(--color-level-2);
					width: 134px;
					height: 18px;
				}
				.playback {
					width: max-content;
					height: max-content;
				}
			}
		}
	}
</style>
