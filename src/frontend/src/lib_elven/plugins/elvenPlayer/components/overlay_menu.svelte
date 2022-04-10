<script lang="ts">
    import { createEventDispatcher } from "svelte";
    // types
    import type { ComponentState } from "../types";
    // ui
    import Overlay from "./overlay.svelte";
    import Progress from "./progress.svelte";
    import Slider from "./slider.svelte";
    import SliderTime from "./slider_time.svelte";

    export let state: ComponentState;

    const dispatch = createEventDispatcher<{
        /** in percents */
        volumeChanged: number;
        /** in percents */
        currentTimeChanged: number;
        /** when user drag time slider. In percents */
        currentTimePreviewChanged: number;
        /** overlay on/off */
        deactivated: void;
    }>();

    /** current player state */

    function onVolumeChanged(perc: number) {
        dispatch("volumeChanged", perc);
    }

    function onCurrentTimeChanged(perc: number) {
        dispatch("currentTimeChanged", perc);
    }

    function onCurrentTimePreviewChanged(perc: number) {
        dispatch("currentTimePreviewChanged", perc);
    }

    function setCurrentTimeDraggingNow(v: boolean) {
        state.current.time.draggingNow = v;
    }

    function onDeactivated() {
        dispatch("deactivated");
    }
</script>

<Overlay on:deactivated={() => onDeactivated()}>
    <div class="overlay">
        <div class="current">
            <div class="current__sliders">
                <div class="progress__buffered">
                    <Progress
                        bind:percents={state.current.buffered.percents}
                        --color="#383659"
                        --border-radius="12px"
                    />
                </div>
                <div class="slider__time">
                    <SliderTime
                        positionPercents={state.current.time.percents}
                        on:draggingNow={(e) =>
                            setCurrentTimeDraggingNow(e.detail)}
                        on:currentTimeChanged={(e) =>
                            onCurrentTimeChanged(e.detail)}
                        on:currentTimePreview={(e) =>
                            onCurrentTimePreviewChanged(e.detail)}
                    />
                </div>
            </div>
            <div class="current__info">
                <div class="current__position">
                    {state.current.time.pretty}
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
            <slot name="playbackControls" />
        </div>
    </div>
</Overlay>

<style lang="scss">
    .overlay {
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
            background-color: #000;
            border-radius: 4px;
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
    }

    .slider__volume {
        background-color: #000;
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
