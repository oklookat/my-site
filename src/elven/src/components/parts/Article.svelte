<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { TArticle } from "@/types/article";
    import Dates from "@/tools/Dates";

    export let article: TArticle;

    const dispatch = createEventDispatcher<{ selected: TArticle }>();

    function convertDate(date) {
        return Dates.convert(date);
    }

    function onSelected(article: TArticle) {
        dispatch("selected", article);
    }
</script>

<article class="article" on:click={() => onSelected(article)}>
    <div class="article__meta">
        {#if !article.is_published}
            <div class="article__item article__updated">
                {convertDate(article.updated_at)}
            </div>
        {:else}
            <div class="article__item article__published">
                {convertDate(article.published_at)}
            </div>
        {/if}
        {#if !article.is_published}
            <div class="article__item article__is-draft">draft</div>
        {/if}
    </div>
    <div class="article__main">
        <div class="article__item article__title">
            {article.title}
        </div>
        {#if article.content && article.content.blocks && article.content.blocks[0]}
            <div class="article__item article__preview">
                <div>{@html article.content.blocks[0].data.text}</div>
            </div>
        {/if}
    </div>
</article>

<style lang="scss">
    .article {
        min-height: 164px;
        border-radius: var(--border-radius);
        background-color: var(--color-level-1);
        cursor: pointer;
        width: 100%;
        display: flex;
        flex-direction: column;
        padding-bottom: 12px;
        gap: 8px;
        &__item {
            font-size: 1.1rem;
            line-height: 1.5rem;
            margin-top: 8px;
            margin-left: 12px;
            margin-right: 12px;
        }
        &__meta {
            display: flex;
            flex-direction: row;
            color: var(--color-text-inactive);
        }
        &__main {
            display: flex;
            flex-direction: column;
            gap: 4px;
        }
        &__title {
            font-size: 1.5rem;
            line-height: 2rem;
            letter-spacing: 0.0099rem;
        }
        &__preview {
            height: 100%;
        }
        &__updated,
        &__published,
        &__is-draft {
            font-size: 0.9rem;
        }
    }
</style>
