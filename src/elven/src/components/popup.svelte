<script lang="ts">
    import { onMount } from "svelte";

    export let mouseEvent: MouseEvent;

    export let onDisabled: () => void;

    let popupEL: HTMLDivElement;
    let initialClick = true;

    onMount(() => {
        enable(mouseEvent);
    });

    function enable(evt: MouseEvent) {
        if (!evt) {
            return;
        }
        const { x, y } = correctOverflow(evt);
        // set styles
        popupEL.style.left = `${x}px`;
        popupEL.style.top = `${y}px`;
    }

    /** check is popup not out of screen, and if it is, correct position */
    function correctOverflow(evt: MouseEvent): { x: number; y: number } {
        let x = evt.clientX;
        let y = evt.clientY;
        const moveOffset = 10;

        // left-right (X)
        const popupWidth = popupEL.offsetWidth;
        const overflowDifferenceX = x + popupWidth - window.screen.width;
        if (overflowDifferenceX > 0) {
            x = x - overflowDifferenceX - moveOffset;
        }

        // top-bottom (Y)
        const popupHeight = popupEL.offsetHeight;
        const overflowDifferenceY = y + popupHeight - window.screen.height;
        if (overflowDifferenceY > 0) {
            y = y - overflowDifferenceY - moveOffset;
        }

        return { x, y };
    }

    function watchClick(e: MouseEvent) {
        if (initialClick) {
            initialClick = false;
            return;
        }
        if (!onDisabled) {
            return;
        }
        e.preventDefault();
        const target = e.target;
        if (target instanceof HTMLElement) {
            const isChild = isChildOfPopup(target);
            if (isChild) {
                return;
            }
        }
        onDisabled();
    }

    function isChildOfPopup(el: HTMLElement): boolean {
        if (el === popupEL) {
            return true;
        }
        while (
            el.parentNode &&
            el.parentNode.nodeName.toLowerCase() !== "body"
        ) {
            el = el.parentNode as any;
            if (el === popupEL) {
                return true;
            }
        }
        return false;
    }
</script>

<svelte:window on:click={watchClick} />
<div class="popup with-border" bind:this={popupEL}>
    <div class="item">итем</div>
    <div class="item">итем</div>
    <div class="item">итем</div>
    <div class="item">итем</div>
    <div class="item">итем</div>
</div>

<style lang="scss">
    .popup {
        display: block;
        position: absolute;
        background-color: var(--color-level-1);
        border-radius: var(--border-radius);
        width: 200px;
        height: max-content;
        padding: 6px;
        display: flex;
        flex-direction: column;
        gap: 12px;
        :global(div),
        :global(a) {
            border-bottom: 1px solid var(--color-body);
            width: 100%;
            height: auto;
            min-height: 32px;
            &:hover {
                background-color: var(--color-hover);
            }
        }
    }
</style>
