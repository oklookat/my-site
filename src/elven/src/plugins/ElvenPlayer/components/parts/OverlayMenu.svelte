<script lang="ts">
    import { onDestroy } from "svelte";
    import type Core from "../../core";
    import Overlay from "../ui/Overlay.svelte";
    import PlaybackControls from "./PlaybackControls.svelte";
    import Progress from "../ui/Progress.svelte";
    import Slider from "../ui/slider/Slider.svelte";
    import type { Unsubscriber } from "svelte/store";
    import TimeSlider from "./TimeSlider.svelte";

    export let core: Core;
    export let active: boolean;

    /** current player state */
    const state = {
        volume: {
            percents: 0,
        },
        current: {
            buffered: {
                percents: 0,
            },
            position: {
                draggingNow: false,
                percents: 0,
                pretty: "00:00",
            },
            duration: {
                pretty: "00:00",
            },
        },
    };

    const unsubs: Unsubscriber[] = [];

    const u1 = core.state.store.current.buffered.percents.subscribe((v) => {
        state.current.buffered.percents = v;
    });
    unsubs.push(u1);

    const u2 = core.state.store.current.position.percents.subscribe((v) => {
        state.current.position.percents = v;
    });
    unsubs.push(u2);

    const u3 = core.state.store.current.position.pretty.subscribe((v) => {
        // not setting pretty if user dragging time slider now (time preview)
        if (state.current.position.draggingNow) {
            return;
        }
        state.current.position.pretty = v;
    });
    unsubs.push(u3);

    const u4 = core.state.store.current.duration.pretty.subscribe((v) => {
        state.current.duration.pretty = v;
    });
    unsubs.push(u4);

    const u5 = core.state.store.volume.percents.subscribe((v) => {
        state.volume.percents = v;
    });
    unsubs.push(u5);

    function onVolumeChanged(perc: number) {
        core.dom.volumePercents = perc;
    }

    function setPositionPercents(perc: number) {
        core.dom.currentTimePercents = perc;
    }

    function setPositionPreview(perc: number) {
        state.current.position.pretty =
            core.dom.getCurrentTimePrettyByPercents(perc);
    }

    function setPositionDraggingNow(v: boolean) {
        state.current.position.draggingNow = v;
    }

    onDestroy(() => {
        for (const unsub of unsubs) {
            unsub();
        }
    });
</script>

{#if active}
    <Overlay on:deactivated={() => (active = false)}>
        <div class="overlay__controls">
            <div class="current">
                <div class="current__sliders">
                    <div class="progress__buffered">
                        <Progress
                            bind:percents={state.current.buffered.percents}
                        />
                    </div>
                    <div class="slider__time">
                        <TimeSlider
                            positionPercents={state.current.position.percents}
                            on:positionSet={(e) =>
                                setPositionPercents(e.detail)}
                            on:positionPreview={(e) =>
                                setPositionPreview(e.detail)}
                            on:draggingNow={(e) =>
                                setPositionDraggingNow(e.detail)}
                        />
                    </div>
                </div>
                <div class="current__info">
                    <div class="current__position">
                        {state.current.position.pretty}
                    </div>
                    <div class="current__total">
                        {state.current.duration.pretty}
                    </div>
                </div>
            </div>

            <div class="slider__volume">
                <Slider
                    percents={state.volume.percents}
                    afterUp={false}
                    on:slide={(e) => onVolumeChanged(e.detail)}
                />
            </div>

            <div class="playback">
                <PlaybackControls {core} />
            </div>
        </div>
    </Overlay>
{/if}

<style lang="scss">
    .overlay__controls {
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
