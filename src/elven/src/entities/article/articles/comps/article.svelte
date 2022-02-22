<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { Article } from "../types";
    import Dates from "@/tools/dates";

    export let article: Article;

    const dispatch = createEventDispatcher<{ selected: Article }>();

    function convertDate(date: string | number | Date): string {
        return Dates.convert(date);
    }

    function onSelected(article: Article) {
        dispatch("selected", article);
    }
</script>

<article class="article" on:click={() => onSelected(article)}>
    <div class="article__meta">
        <div class="article__item meta__item">
            {#if article.is_published && article.published_at}
                {convertDate(article.published_at)}
            {:else}
                {convertDate(article.updated_at)}
            {/if}
        </div>
        <div class="article__item meta__item">
            {#if !article.is_published}
                draft
            {:else}
                published
            {/if}
        </div>
        {#if article.category_name}
            <div class="article__item meta__item">
                {article.category_name}
            </div>
        {/if}
    </div>
    <div class="article__main">
        <div class="article__item article__title">
            {article.title}
        </div>
    </div>
</article>

<style lang="scss">
    .article {
        border-radius: var(--border-radius);
        background-color: var(--color-level-1);
        cursor: pointer;
        width: 100%;
        height: max-content;
        display: flex;
        flex-direction: column;
        &__item {
            font-size: 1.1rem;
        }
        &__meta {
            padding: 12px 12px 8px;
            display: flex;
            flex-direction: row;
            height: fit-content;
            gap: 12px;
            flex-wrap: wrap;
            .meta__item {
                font-size: 0.9rem;
                min-width: 44px;
                border-radius: 12px;
                min-height: 32px;
                height: fit-content;
                display: flex;
                justify-content: center;
                align-items: center;
            }
        }
        &__main {
            width: 100%;
            height: 100%;
            display: flex;
            flex-direction: column;
            gap: 4px;
        }
        &__title {
            font-size: 1.6rem;
            line-height: 2rem;
            letter-spacing: 0.0099rem;
            padding: 8px 12px 18px;
        }
    }
</style>
