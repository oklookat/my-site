<template>
    <div class="articles__list" v-if="articles.length > 0">
        <article
            class="article"
            v-for="article in articles"
            :key="article.id"
            v-on:click="onSelected(article)"
        >
            <div class="article__meta">
                <div
                    class="article__item article__updated"
                    v-if="!article.is_published"
                >{{ convertDate(article.updated_at) }}</div>
                <div
                    class="article__item article__published"
                    v-else
                >{{ convertDate(article.published_at) }}</div>
                <div class="article__item article__is-draft" v-if="!article.is_published">draft</div>
            </div>
            <div class="article__main">
                <div class="article__item article__title">{{ article.title }}</div>
                <div
                    v-if="article.content && article.content.blocks && article.content.blocks[0]"
                    class="article__item article__preview"
                >
                    <div v-html="article.content.blocks[0].data.text"></div>
                </div>
            </div>
        </article>
    </div>
</template>

<script setup lang="ts">
import { IArticle } from '@/types/article'
import Dates from '@/common/tools/Dates'

const props = defineProps<{
    articles: Array<IArticle>
}>()


const emit = defineEmits<{
    (e: 'selected', article: IArticle): void
}>()

function convertDate(date) {
    return Dates.convert(date)
}

function onSelected(article: IArticle) {
    emit('selected', article)
}
</script>

<style scoped>
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