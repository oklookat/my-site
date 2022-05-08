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
	import { onDestroy, onMount } from 'svelte';
	import PlaybackControls from './playback_controls.svelte';
	import { toggleBodyScroll } from '$lib/tools';

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

	let defScroll: () => void;
	onMount(() => {
		defScroll = toggleBodyScroll();
	});

	onDestroy(() => {
		defScroll();
		unsub1();
	});

	let slidedPercents = 0;
	function onProgressSlide(percents: number) {
		slidedPercents = percents;
		currentTimeNew = convertPercentsToCurrentTime(slidedPercents, $duration);
		currentTimePretty = getPretty(currentTimeNew);
	}
</script>

<div class="overlay" on:click|self|stopPropagation={onClose}>
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

			<div class="volume">
				<div class="itself">
					<Slider percents={$volumePercents} on:slide={(e) => ($volumePercents = e.detail)} />
				</div>
			</div>

			<div class="controls">
				<PlaybackControls />
			</div>
		</div>
	</div>
</div>

<style lang="scss">
	.overlay {
		background-color: rgba(0, 0, 0, 0.4);
		z-index: 9998;
		max-width: 100vw;
		width: 100%;
		height: 100%;
		position: fixed;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		.main {
			user-select: none;
			background-color: var(--color-level-1);
			border-radius: var(--border-radius);
			height: 314px;
			width: 424px;
			@media screen and(max-width: 644px) {
				width: 95%;
			}
			padding: 24px;

			> div {
				width: 100%;
				height: 100%;
				display: flex;
				flex-direction: column;

				div {
					height: 50%;
				}

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

				.volume,
				.controls {
					display: flex;
					align-items: center;
					justify-content: center;
				}
				.volume {
					height: 12px;
					width: 100%;
					.itself {
						height: 100%;
						width: 112px;
					}
				}
			}
		}
	}
</style>
