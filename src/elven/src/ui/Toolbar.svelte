<div class="toolbar">
    <div class="toolbar__content">
        <slot />
    </div>
</div>

<style lang="scss">
    .toolbar {
        font-size: 0.9rem;
        background-color: var(--color-level-1);
        border-radius: var(--border-radius);
        width: 100%;
        min-height: 52px;
        padding: 14px;
        &__content {
            height: 100%;
            display: flex;
            flex-flow: row wrap;
            align-items: center;
            align-content: center;
            gap: 8px;
        }
    }
</style>
