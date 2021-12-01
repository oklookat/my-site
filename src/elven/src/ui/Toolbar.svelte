<div class="toolbar">
    <div class="toolbar__content">
        <slot />
    </div>
</div>

<style lang="scss">
    $padding: 14px;

    .toolbar {
        background-color: var(--color-level-1);
        border-radius: var(--border-radius);
        font-size: 0.8rem;
        height: 36px;
        width: 100%;
        &__content {
            padding-left: $padding;
            padding-right: $padding;
            display: flex;
            flex-flow: row wrap;
            gap: 8px;
            height: 100%;
            align-items: center;
        }
    }
</style>
