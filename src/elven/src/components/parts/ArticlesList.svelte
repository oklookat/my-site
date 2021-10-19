<script lang="ts">
    import type { IArticle } from "@/types/article";
    import Dates from "@/common/tools/Dates";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let articles: Array<IArticle>;

    function convertDate(date) {
        return Dates.convert(date);
    }

    function onSelected(article: IArticle) {
        dispatch("selected", article);
    }
</script>

<main>
    {#if articles.length > 0}
        <div class="articles__list">
            {#each articles as article (article.id)}
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
                            <div class="article__item article__is-draft">
                                draft
                            </div>
                        {/if}
                    </div>
                    <div class="article__main">
                        <div class="article__item article__title">
                            {article.title}
                        </div>
                        {#if article.content && article.content.blocks && article.content.blocks[0]}
                            <div class="article__item article__preview">
                                <div
                                    v-html="article.content.blocks[0].data.text"
                                />
                            </div>
                        {/if}
                    </div>
                </article>
            {/each}
        </div>
    {/if}
</main>

<style>
    .articles__list {
        height: 100%;
        width: 100%;
        display: flex;
        flex-direction: column;
        min-height: 42px;
        gap: 12px;
    }

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
    }

    .article__item {
        font-size: 1.1rem;
        line-height: 1.5rem;
        margin-top: 8px;
        margin-left: 12px;
        margin-right: 12px;
    }

    .article__meta {
        display: flex;
        flex-direction: row;
        color: var(--color-text-inactive);
    }

    .article__main {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .article__title {
        font-size: 1.5rem;
        line-height: 2rem;
        letter-spacing: 0.0099rem;
    }

    .article__preview {
        height: 100%;
    }

    .article__updated,
    .article__published,
    .article__is-draft {
        font-size: 0.9rem;
    }
</style>
