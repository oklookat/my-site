<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { Article } from "../types";
    import Dates from "@/tools/dates";

    export let article: Article;

    const dispatch = createEventDispatcher<{ selected: Article }>();

    function convertDate(date: string | number | Date) {
        return Dates.convert(date);
    }

    function onSelected(article: Article) {
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
    </div>
</article>

<style lang="scss">
    .article {
        min-height: 94px;
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
        &__updated,
        &__published,
        &__is-draft {
            font-size: 0.9rem;
        }
    }
</style>
