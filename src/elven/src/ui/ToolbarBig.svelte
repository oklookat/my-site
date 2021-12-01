<div class="toolbar">
    <div class="toolbar__content">
        <slot />
    </div>
</div>

<style lang="scss">
    $padding: 14px;

    .toolbar {
        border-radius: var(--border-radius);
        background-color: var(--color-level-1);
        height: 48px;
        width: 100%;
        font-size: 1.1rem;
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
