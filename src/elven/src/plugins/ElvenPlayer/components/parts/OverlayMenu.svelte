<script lang="ts">
    import { onDestroy } from "svelte";
    import type Core from "../../core";
    import Overlay from "../ui/Overlay.svelte";
    import PlaybackControls from "./PlaybackControls.svelte";
    import Progress from "../ui/Progress.svelte";
    import Slider from "../ui/slider/Slider.svelte";

    export let core: Core;
    export let active: boolean;

    function onProgressSliderTriggered(percents: number) {
        core.dom.currentTimePercents = percents;
    }

    function onVolumeSliderTriggered(percents: number) {
        core.dom.volumePercents = percents;
    }

    // state
    let bufferedPercents = 0;
    const s1 = core.state.store.current.buffered.percents.subscribe((v) => {
        bufferedPercents = v;
    });

    let currentTimePretty = "00:00";
    const s2 = core.state.store.current.position.pretty.subscribe((v) => {
        currentTimePretty = v;
    });

    let currentTimePercents = 0;
    const s3 = core.state.store.current.position.percents.subscribe((v) => {
        currentTimePercents = v;
    });

    let durationPretty = "00:00";
    const s4 = core.state.store.current.duration.pretty.subscribe((v) => {
        durationPretty = v;
    });

    let volumePercents = 100;
    const s5 = core.state.store.volume.percents.subscribe((v) => {
        volumePercents = v;
    });

    onDestroy(() => {
        s1();
        s2();
        s3();
        s4();
        s5();
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
                            on:slide={(e) => onProgressSliderTriggered(e.detail)}
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
