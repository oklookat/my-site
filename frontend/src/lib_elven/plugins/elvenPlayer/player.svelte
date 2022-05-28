<script lang="ts">
	import { browser } from '$app/env';
	import Close from './icons/close.svelte';
	import Audio from './audio.svelte';
	import { audioComponent, currentSource, isPlaying, setCurrentTime, signal } from './store';
	import { Signal, type ElvenPlayer, type Playlist } from './types';
	import { onDestroy, onMount } from 'svelte';
	import Hamburger from './icons/hamburger.svelte';
	import PlaybackControls from './playback_controls.svelte';
	import PlaybackMenu from './playback_menu.svelte';

	let isActive = false;

	/** current playlist */
	const playlist: Playlist = {
		currentPosition: 0,
		position: 0,
		sources: []
	};

	class Plugin implements ElvenPlayer {
		addToPlaylist(src: URL) {
			for (const source of playlist.sources) {
				if (src.toString() === source.toString()) {
					return;
				}
			}
			playlist.sources.push(src);
		}

		clearPlaylist() {
			playlist.position = 0;
			playlist.currentPosition = 0;
			playlist.sources = [];
		}

		async play() {
			if (!isActive) {
				isActive = true;
			}

			playlist.currentPosition = playlist.position;
			const newSrc = playlist.sources[playlist.currentPosition];

			if (!newSrc) {
				console.warn('[elvenPlayer] no sources in playlist.');
				return;
			}

			const isSameSource = !!(
				$audioComponent &&
				$currentSource &&
				$currentSource.toString() === newSrc.toString()
			);
			if (isSameSource) {
				if ($isPlaying) {
					this.pause();
				} else {
					$signal = Signal.PLAY;
				}
				return;
			}

			await this.playCurrentSource();
		}

		pause() {
			$signal = Signal.PAUSE;
		}

		async next() {
			const newSrc = playlist.sources[playlist.position + 1];
			if (newSrc) {
				playlist.position++;
				await this.play();
				return;
			}
			$setCurrentTime = 0;
		}

		async prev() {
			const prevSrc = playlist.sources[playlist.position - 1];
			if (prevSrc) {
				playlist.position--;
				await this.play();
			}

			await this.playCurrentSource();
		}

		private async playCurrentSource() {
			await createSource();
			$signal = Signal.PLAY;
		}
	}

	onMount(() => {
		if (!browser) {
			return;
		}
		console.log('MOUNT');
		window.$player = new Plugin();
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		console.log('DESTROY');
		close();
		window.$player = undefined;
	});

	let container: HTMLDivElement;

	/** (re)create source (component) with current source (URL) */
	async function createSource(): Promise<void> {
		destroyAudio();
		/** why use setTimeout 0 and promise? One of the cases: */
		// look at the play() function
		// we have isActive check. If isActive === false -> we set to true ->
		// "player" class would be mounted -> we provide this class as target to audioComponent -> play audio
		// BUT due on svelte-specific (or not) as i understand DOM conditional rendering
		// executes in the end of order.
		// that is, despite the fact that createSource() is called later, it will be called earlier than
		// "player" class would be rendered -> and we can't provide this class as target.
		// But if we use setTimeout 0 we can put audioComponent creating to the end of the order
		// and when promise resolves we know - "player" class ready.
		return new Promise((resolve) => {
			setTimeout(() => {
				$audioComponent = new Audio({
					target: container,
					props: {
						onEnded: () => {
							window.$player?.next();
						},
						source: playlist.sources[playlist.currentPosition]
					}
				});
				resolve();
			}, 0);
		});
	}

	function destroyAudio() {
		$audioComponent?.$destroy();
		audioComponent.set(undefined);
	}

	/** close player UI and destroy audio */
	function close() {
		window.$player?.clearPlaylist();
		destroyAudio();
		isActive = false;
	}

	let isPlaybackMenuActive = false;
</script>

{#if isPlaybackMenuActive}
	<PlaybackMenu onClose={() => (isPlaybackMenuActive = false)} />
{/if}

{#if isActive}
	<div class="player" bind:this={container}>
		<div>
			<div class="menu" on:click={() => (isPlaybackMenuActive = !isPlaybackMenuActive)}>
				<Hamburger />
			</div>
			<div class="control">
				<PlaybackControls />
			</div>
			<div class="close" on:click={close}>
				<Close />
			</div>
		</div>
	</div>
{/if}

<style lang="scss">
	@import './icons/vars.scss';

	.player {
		width: 100%;
		height: 64px;
		background-color: var(--color-level-1);

		> div {
			height: 100%;
			width: 95%;
			margin: auto;
			div {
				width: max-content;
				height: 100%;
				display: flex;
				align-items: center;
				justify-content: center;
			}

			display: grid;
			grid-template-columns: repeat(3, 1fr);
			justify-content: center;
			align-items: center;

			.control {
				width: 100%;
			}

			.close {
				justify-self: flex-end;
			}
		}
	}
</style>
