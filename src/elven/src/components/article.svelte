<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { Article } from "@/types/articles";
    import Dates from "@/tools/dates";
    import ArticleCover from "./article_cover.svelte";

    export let article: Article;

    const dispatch = createEventDispatcher<{ selected: Article }>();

    function convertDate(date: string | number | Date): string {
        return Dates.convert(date);
    }

    function onSelected(article: Article) {
        dispatch("selected", article);
    }
</script>

<article class="article base__card" on:click={() => onSelected(article)}>
    <div class="meta">
        <div class="meta__item">
            {#if article.is_published && article.published_at}
                {convertDate(article.published_at)}
            {:else}
                {convertDate(article.updated_at)}
            {/if}
        </div>
        <div class="meta__item">
            {#if !article.is_published}
                draft
            {:else}
                published
            {/if}
        </div>
        {#if article.category_name}
            <div class="meta__item">
                {article.category_name}
            </div>
        {/if}
    </div>
    <div class="main">
        {#if article.cover_id && article.cover_path && article.cover_extension}
            <div class="cover">
                <ArticleCover {article} />
            </div>
        {/if}
        <div class="title">
            {article.title}
        </div>
    </div>
</article>

