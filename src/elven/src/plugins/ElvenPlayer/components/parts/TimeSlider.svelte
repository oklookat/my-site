<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import Slider from "../ui/slider/Slider.svelte";

    const dispatch = createEventDispatcher<{
        /** when user dragging slider or not */
        draggingNow: boolean;
        /** when user dragging slider returns percents */
        currentTimePreview: number;
        /** after dragging slider returns percents */
        currentTimeChanged: number;
    }>();

    /** is user currently dragging slider */
    let draggingNow = false;

    /** temp position percents for preview. Setting to positionPercents after mouse up. */
    let tempPositionPercents = 0;

    /** current audio position in percents */
    export let positionPercents: number;

    function onSlide(perc: number) {
        // if mouse up after slide time, we set audio position by percents
        if (!draggingNow) {
            dispatch("currentTimeChanged", perc);
            return;
        }
        // if mouse down, we write time percents to buff, and calculate pretty preview
        tempPositionPercents = perc;
        dispatch("currentTimePreview", perc);
    }

    function onMouse(v: boolean) {
        draggingNow = v;
        dispatch("draggingNow", draggingNow);
        // if not dragging set position
        if (!draggingNow) {
            dispatch("currentTimeChanged", tempPositionPercents);
        }
    }
</script>

<Slider
    afterUp={false}
    on:slide={(e) => onSlide(e.detail)}
    bind:percents={positionPercents}
    on:mouse={(e) => onMouse(e.detail)}
/>
