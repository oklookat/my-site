<script lang="ts">
    import { onDestroy } from "svelte";
    import type Core from "../../core";
    import Overlay from "../ui/Overlay.svelte";
    import PlaybackControls from "./PlaybackControls.svelte";
    import Progress from "../ui/Progress.svelte";
    import Slider from "../ui/slider/Slider.svelte";
    import type { Unsubscriber } from "svelte/store";

    export let core: Core;
    export let active: boolean;

    const onProgressSliderTriggered = (perc: number) =>
        (core.dom.currentTimePercents = perc);
    const onVolumeSliderTriggered = (perc: number) =>
        (core.dom.volumePercents = perc);

    const unsubs: Unsubscriber[] = [];

    let bufferedPercents: number;
    const u1 = core.state.store.current.buffered.percents.subscribe((v) => {
        bufferedPercents = v;
    });
    unsubs.push(u1);

    let currentTimePretty: string;
    const u2 = core.state.store.current.position.pretty.subscribe((v) => {
        currentTimePretty = v;
    });
    unsubs.push(u2);

    let currentTimePercents: number;
    const u3 = core.state.store.current.position.percents.subscribe((v) => {
        currentTimePercents = v;
    });
    unsubs.push(u3);

    let durationPretty: string;
    const u4 = core.state.store.current.duration.pretty.subscribe((v) => {
        durationPretty = v;
    });
    unsubs.push(u4);

    let volumePercents: number;
    const u5 = core.state.store.volume.percents.subscribe((v) => {
        volumePercents = v;
    });
    unsubs.push(u5);

    onDestroy(() => {
        for (const unsub of unsubs) {
            unsub();
        }
    });
</script>

{#if active}
    <Overlay on:deactivated={() => (active = false)}>
        <div class="overlay__menu">
            <div class="current">
                <div class="current__sliders">
                    <div class="progress__buffered">
                        <Progress bind:percents={bufferedPercents} />
                    </div>
                    <div class="slider__time">
                        <Slider
                            afterUp={true}
                            on:slide={(e) =>
                                onProgressSliderTriggered(e.detail)}
                            bind:percents={currentTimePercents}
                        />
                    </div>
                </div>
                <div class="current__info">
                    <div class="current__position">
                        {currentTimePretty}
                    </div>
                    <div class="current__total">
                        {durationPretty}
                    </div>
                </div>
            </div>

            <div class="slider__volume">
                <Slider
                    percents={volumePercents}
                    afterUp={false}
                    on:slide={(e) => onVolumeSliderTriggered(e.detail)}
                />
            </div>

            <div class="playback">
                <PlaybackControls {core} />
            </div>
        </div>
    </Overlay>
{/if}

<style lang="scss">
    .overlay__menu {
        height: 100%;
        width: 100%;
        display: grid;
        grid-template-columns: 1fr;
        align-items: center;
        justify-items: center;
    }
    .current {
        width: 85%;
        display: grid;
        grid-template-columns: 1fr;
        grid-template-rows: 1fr 1fr;
        gap: 18px;
        &__sliders {
            position: relative;
            .progress__buffered,
            .slider__time {
                position: absolute;
                width: 100%;
                height: 100%;
            }
        }
        &__info {
            display: flex;
            flex-direction: row;
            .current__total {
                margin-left: auto;
            }
        }
    }

    .slider__time,
    .slider__volume,
    .progress__buffered {
        border-radius: 4px;
        background-color: rgba(0, 0, 0, 0.5);
    }

    .slider__volume {
        position: relative;
        width: 50%;
        height: 14px;
    }

    .playback {
        width: 100%;
        height: 100%;
        display: flex;
        justify-content: center;
        fill: black;
        @media (prefers-color-scheme: dark) {
            fill: white;
        }
    }
</style>
