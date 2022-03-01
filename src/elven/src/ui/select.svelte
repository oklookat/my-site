<script lang="ts">
    import { createEventDispatcher } from "svelte";

    /** default selected value */
    export let selected: string = undefined;
    /** options */
    export let options: {
        value: string;
        text: string;
    }[] = [];

    const dispatch = createEventDispatcher<{
        /** when other item selected */
        selected: string;
    }>();

    function onChange(e: Event) {
        const target = e.target as HTMLSelectElement;
        dispatch("selected", target.value);
    }
</script>

<select class="select" on:change={(e) => onChange(e)} value={selected}>
    {#each options as piece}
        <option value={piece.value}>{piece.text}</option>
    {/each}
</select>

<style lang="scss">
    .select {
        color: var(--color-text);
        background-color: var(--color-level-2);
        box-sizing: border-box;
        border-radius: 0.4rem;
        min-width: 94px;
        max-width: fit-content;
        height: 100%;
        padding: 4px;
    }
</style>
