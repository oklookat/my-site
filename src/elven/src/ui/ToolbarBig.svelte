<div class="toolbar">
    <div class="toolbar__content">
        <slot />
    </div>
</div>

<style lang="scss">
    .toolbar {
        font-size: 1.1rem;
        border-radius: var(--border-radius);
        background-color: var(--color-level-1);
        width: 100%;
        min-height: 42px;
        &__content {
            padding: 12px;
            min-height: 42px + 12px;
            height: 100%;
            display: flex;
            flex-flow: row wrap;
            align-items: center;
            gap: 14px;
        }
    }
</style>
