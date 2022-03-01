<div class="toolbar with-border">
    <div class="toolbar__content">
        <slot />
    </div>
</div>

<style lang="scss">
    .toolbar {
        font-size: 1.1rem;
        background-color: var(--color-level-1);
        width: 100%;
        min-height: 52px;
        &__content {
            padding: 12px;
            height: 100%;
            display: flex;
            flex-flow: row wrap;
            align-items: center;
            align-content: center;
            gap: 14px;
        }
    }
</style>
