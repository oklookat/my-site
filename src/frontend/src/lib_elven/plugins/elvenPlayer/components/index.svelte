<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import type { Unsubscriber } from 'svelte/store';
	// types
	import type { ElvenPlayer, ComponentState } from '../types';
	// ui
	import PlaybackControls from './playback_controls.svelte';
	import OverlayMenu from './overlay_menu.svelte';
	import Hamburger from '$lib/icons/hamburger.svelte';

	/** is player active */
	let active: boolean = false;
	export let core: ElvenPlayer;

	/** when player open/close */
	export let onActiveChanged: ((active?: boolean) => void) | undefined = undefined;
	$: onActiveChange(active);
	function onActiveChange(val: boolean) {
		if (!onActiveChanged) {
			return;
		}
		onActiveChanged(val);
	}

	/** is player controls overlay active */
	let isOverlay = false;

	/** player state */
	let state: ComponentState = {
		playing: false,
		volume: {
			percents: 100
		},
		current: {
			buffered: {
				percents: 0
			},
			time: {
				draggingNow: false,
				percents: 0,
				pretty: '00:00'
			},
			duration: {
				pretty: '00:00'
			}
		}
	};

	let unsubs: Unsubscriber[] = [];

	onMount(() => {
		init();
	});

	onDestroy(() => {
		destroy();
	});

	function init() {
		unsubs.push(
			core.store.state.playing.onChange((v) => {
				if (typeof v !== 'boolean') {
					return;
				}
				if (v) {
					active = true;
				}
				state.playing = v;
			})
		);

		unsubs.push(
			core.store.state.current.buffered.percents.onChange((v) => {
				if (typeof v !== 'number') {
					return;
				}
				state.current.buffered.percents = v;
			})
		);

		unsubs.push(
			core.store.state.current.time.percents.onChange((v) => {
				if (typeof v !== 'number') {
					return;
				}
				state.current.time.percents = v;
			})
		);

		unsubs.push(
			core.store.state.current.time.pretty.onChange((v) => {
				// not setting pretty if user dragging time slider now (time preview)
				if (state.current.time.draggingNow) {
					return;
				}
				if (typeof v !== 'string') {
					return;
				}
				state.current.time.pretty = v;
			})
		);

		unsubs.push(
			core.store.state.current.duration.pretty.onChange((v) => {
				if (typeof v !== 'string') {
					return;
				}
				state.current.duration.pretty = v;
			})
		);

		unsubs.push(
			core.store.state.volume.percents.onChange((v) => {
				if (typeof v !== 'number') {
					return;
				}
				state.volume.percents = v;
			})
		);
	}

	function destroy() {
		for (const unsub of unsubs) {
			unsub();
		}
		unsubs = [];
		active = false;
		core.stop();
	}

	function onPlayPause(isPlay: boolean) {
		active = true;
		core.playPause()
	}

	function onNext() {
		core.next();
	}

	function onPrev() {
		core.prev();
	}

	function onClose() {
		active = false;
		core.stop();
	}

	function setVolumePercents(perc: number) {
		core.volumePercents = perc;
	}

	function setCurrentTimePercents(perc: number) {
		core.currentTimePercents = perc;
	}

	function setCurrentTimePreview(perc: number) {
		state.current.time.pretty = core.convertPercentsToCurrentTimePretty(perc);
	}
</script>

{#if isOverlay}
	<OverlayMenu
		bind:state
		on:deactivated={() => (isOverlay = false)}
		on:volumeChanged={(e) => setVolumePercents(e.detail)}
		on:currentTimeChanged={(e) => setCurrentTimePercents(e.detail)}
		on:currentTimePreviewChanged={(e) => setCurrentTimePreview(e.detail)}
	>
		<PlaybackControls
			slot="playbackControls"
			bind:isPlaying={state.playing}
			on:playPause={e => (onPlayPause(e.detail))}
			on:next={() => onNext()}
			on:prev={() => onPrev()}
		/>
	</OverlayMenu>
{/if}

{#if active}
	<div class="player">
		<div class="show">
			<Hamburger on:click={() => (isOverlay = !isOverlay)} />
		</div>

		<div class="controls" on:click|stopPropagation>
			<PlaybackControls
				bind:isPlaying={state.playing}
				on:playPause={e => (onPlayPause(e.detail))}
				on:next={() => onNext()}
				on:prev={() => onPrev()}
			/>
		</div>

		<div class="close" on:click={() => onClose()}>
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 255.07 295.91">
				<path
					d="M135,390.18h0a33.69,33.69,0,0,0,48-4.29L369.93,159.34a35.1,35.1,0,0,0-4.19-48.86h0a33.69,33.69,0,0,0-48,4.29L130.83,341.33A35.09,35.09,0,0,0,135,390.18Z"
					transform="translate(-122.85 -102.37)"
				/>
				<path
					d="M365.74,390.18h0a33.68,33.68,0,0,1-48-4.29L130.83,159.34A35.1,35.1,0,0,1,135,110.48h0a33.7,33.7,0,0,1,48,4.29L369.93,341.33A35.09,35.09,0,0,1,365.74,390.18Z"
					transform="translate(-122.85 -102.37)"
				/>
			</svg>
		</div>
	</div>
{/if}

<style lang="scss">
	.player {
		border-top: var(--color-border) 1px solid;
		background-color: var(--color-level-1);
		width: 100%;
		height: 64px;
		display: grid;
		grid-template-columns: 52px 1fr 52px;
		align-items: center;

		div {
			width: 100%;
			height: 50%;
			display: flex;
			justify-content: center;
			align-items: center;
			svg {
				cursor: pointer;
				width: 100%;
				height: 100%;
			}
		}
		
		.close {
			svg {
				height: 60%;
			}
		}

		.controls,
		svg {
			@media (prefers-color-scheme: dark) {
				fill: white;
			}
			@media (prefers-color-scheme: light) {
				fill: black;
			}
		}
	}
</style>
