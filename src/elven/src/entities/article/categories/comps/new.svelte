<script lang="ts">
    import { createEventDispatcher, onDestroy, onMount } from "svelte";

    export let data: string | undefined = undefined;

    const dispatch = createEventDispatcher<{
        submit: string;
        cancel: void;
    }>();
    //
    let textareaEL: HTMLTextAreaElement;
    //
    let _isValid = true;
    /** textarea value valid? if v not provided = getter */
    function isValid(v?: boolean): boolean | void {
        if (v === undefined) {
            return _isValid;
        }
        _isValid = v;
        manageBorder(!_isValid);
    }

    function manageEvents(add: boolean) {
        const action = add ? "addEventListener" : "removeEventListener";
        textareaEL[action]("keydown", onKeydown);
        textareaEL[action]("blur", onBlur);
    }

    /** add or remove red border in textarea */
    function manageBorder(add: boolean) {
        const action = add ? "add" : "remove";
        textareaEL.classList[action]("border-red");
    }

    onMount(() => {
        if (data) {
            textareaEL.value = data;
        }
        textareaEL.focus();
        manageEvents(true);
    });

    onDestroy(() => {
        manageEvents(false);
    });

    /** focus loss */
    function onBlur(e: FocusEvent) {
        dispatch("cancel");
    }

    /** dispatch text on Enter key */
    function onKeydown(e: KeyboardEvent) {
        /**
         * cancel changes if:
         * Escape pressed
         * OR
         * Enter pressed + textarea value empty
         */
        const isEnter = e.key === "Enter";
        const isEsc = e.key === "Escape";
        if (isEnter || isEsc) {
            e.preventDefault();
            const valueEmpty = !textareaEL.value || isSpacesOnly();
            const cancelChanges = isEsc || (isEnter && valueEmpty);
            if (cancelChanges) {
                dispatch("cancel");
                return;
            }
        } else {
            return;
        }
        /**
         * Enter pressed = validate value.
         * If valid - dispatch submit event
         */
        modifyText();
        if (!isValid()) {
            return;
        }
        dispatch("submit", textareaEL.value);
    }

    /** format text before submit */
    function modifyText() {
        if (!textareaEL.value || textareaEL.value.length > 24 || isSpacesOnly()) {
            isValid(false);
            return;
        }
        const modValue = (r: RegExp, to: string) => {
            textareaEL.value = textareaEL.value.replaceAll(r, to);
        };
        // replace new lines with one space
        const newLines = /[\r\n]/gm;
        modValue(newLines, " ");
        // replace 2+ spaces with one space
        const twoOrMoreSpaces = /[^\S]{2,}/gm;
        modValue(twoOrMoreSpaces, " ");
        // remove spaces at start and end
        textareaEL.value = textareaEL.value.trim();
        isValid(true);
    }

    /** value contains spaces only? */
    function isSpacesOnly(): boolean {
        return !textareaEL.value.replace(/\s|\r/g, '').length
    }
</script>

<div class="new">
    <textarea class="new__category" maxlength="24" bind:this={textareaEL} />
</div>

<style lang="scss">
    .new {
        width: 100%;
        &__category {
            overflow: auto;
            width: 100%;
            background-color: transparent;
            color: var(--color-text);
            font-size: 1.3rem;
            padding: 4px;
        }
    }
</style>
