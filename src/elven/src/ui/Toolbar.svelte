<div class="toolbar with-border">
    <div class="toolbar__content">
        <slot />
    </div>
</div>

<style lang="scss">
    .toolbar {
        font-size: 0.9rem;
        background-color: var(--color-level-1);
        width: 100%;
        height: 64px;
        min-height: 64px;
        &__content {
            padding: 12px;
            height: 100%;
            display: flex;
            flex-flow: row wrap;
            align-items: center;
            gap: 14px;
        }
    }
</style>
