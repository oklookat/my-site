<script lang="ts">
	import { browser } from '$app/env';
	import {
		bufferedPercents,
		currentTime,
		currentTimePretty,
		currentTimePercents,
		duration,
		durationPretty,
		isPlaying,
		setCurrentTime,
		volumeFloat,
		volumePercents,
		currentSource,
		signal
	} from './store';

	import { onDestroy } from 'svelte';
	import { getBufferedPercents, getPercents, getPretty } from './utils';
	import { Signal } from './types';

	/** on audio ended */
	export let onEnded: () => void;

	/** audio source*/
	export let source: URL;
	$currentSource = source;

	const unsubSignal = signal.subscribe((newSignal) => {
		switch (newSignal) {
			case Signal.PLAY:
				play();
				return;
			case Signal.PAUSE:
				pause();
				return;
		}
	});

	/** audio element*/
	let audioEL: HTMLAudioElement;

	/** audio element duration */
	let elDuration = 0;

	/** audio element volume */
	let elVolume = 1.0;

	/** is browser and audio element exists? */
	const isReady = () => {
		return !!(browser && audioEL);
	};

	/** change volume on audio element if volume in store changed */
	const unsubVolume = volumePercents.subscribe((val) => {
		if (typeof val !== 'number') {
			return;
		}
		$volumeFloat = val / 100;
		elVolume = $volumeFloat;
		if (!audioEL) {
			return;
		}
		audioEL.volume = elVolume;
	});

	const unsubSetTime = setCurrentTime.subscribe((val) => {
		if (typeof val !== 'number') {
			return;
		}

		if (val >= elDuration || val < 0) {
			val = 0;
		}

		// need this check because
		// 'ended' event can be fired and element destroyed before we set currentTime (?)
		if (!audioEL) {
			return;
		}

		audioEL.currentTime = val;
	});

	onDestroy(() => {
		unsubSignal();
		unsubVolume();
		unsubSetTime();
		if (!isReady()) {
			return;
		}
		pause();
		$currentSource = undefined;
		$isPlaying = false;
		$duration = 0;
		$durationPretty = getPretty(0);
		$currentTime = 0;
		$currentTimePretty = getPretty(0);
		$currentTimePercents = 0;
		$setCurrentTime = 0;
		$bufferedPercents = 0;
	});

	async function play() {
		if (!isReady()) {
			return;
		}
		await audioEL.play();
	}

	function pause() {
		if (!isReady()) {
			return;
		}
		audioEL.pause();
	}

	/** on audio duration changed (usually fires when audio element gets audio) */
	function durationChanged(e: Event) {
		elDuration = audioEL.duration;
		$duration = elDuration;
		$durationPretty = getPretty($duration);
	}

	/** on audio element time updated (when playing) */
	function timeUpdated(e: Event) {
		let elCurrentTime = audioEL.currentTime;
		if (elCurrentTime < 0) {
			elCurrentTime = 0;
		} else if (elCurrentTime > elDuration) {
			elCurrentTime = elDuration;
		}

		$currentTime = elCurrentTime;
		$currentTimePretty = getPretty(elCurrentTime);
		$currentTimePercents = getPercents(elCurrentTime, elDuration);
		$bufferedPercents = getBufferedPercents(elCurrentTime, elDuration, audioEL.buffered);
	}

	/** audioEL.play() called */
	function onPlay(e: Event) {
		$isPlaying = true;
	}

	/** audioEL.pause() called */
	function onPause(e: Event) {
		$isPlaying = false;
	}
</script>

<audio
	crossorigin="anonymous"
	style="display: none;"
	src={source.toString()}
	bind:this={audioEL}
	bind:volume={elVolume}
	on:durationchange={durationChanged}
	on:timeupdate={timeUpdated}
	on:play={onPlay}
	on:pause={onPause}
	on:ended={() => onEnded()}
/>
