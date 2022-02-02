<div class="toolbar">
    <div class="toolbar__content">
        <slot />
    </div>
</div>

<style lang="scss">
    .toolbar {
        font-size: 0.8rem;
        background-color: var(--color-level-1);
        border-radius: var(--border-radius);
        width: 100%;
        min-height: 14px;
        &__content {
            padding: 12px;
            min-height: 12px + 14px;
            height: 100%;
            display: flex;
            flex-flow: row wrap;
            align-items: center;
            gap: 8px;
        }
    }
</style>
