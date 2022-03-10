<script lang="ts">
    // file
    import Extension, {
        type FileTypeSelector,
    } from "@/tools/extension";
    import { PathTools } from "@/tools/paths";
    // article
    import type { Article } from "@/types/articles";

    export let article: Article;
    $: onArticle(article);

    let coverExists = false;
    let extensionType: FileTypeSelector;
    let fullPath: string;
    function onArticle(val: Article) {
        if (!val) {
            return;
        }
        coverExists = !!(
            article.cover_id &&
            article.cover_path &&
            article.cover_extension
        );
        if (!coverExists) {
            return;
        }
        extensionType = Extension.getType(article.cover_extension);
        fullPath = PathTools.getUploadsWith(article.cover_path).toString();
    }
</script>

<div class="cover">
    {#if coverExists}
        {#if extensionType.selected === "IMAGE"}
            <div class="cover__image">
                <img
                    decoding="async"
                    loading="lazy"
                    alt="article cover"
                    src={fullPath}
                />
            </div>
        {:else if extensionType.selected === "VIDEO"}
            <div class="cover__video">
                <video autoplay muted src={fullPath} />
            </div>
        {/if}
    {/if}
</div>

<style lang="scss">
    @import "./src/assets/variables";

    .cover {
        width: 100%;
        height: 100%;
        &__image,
        &__video {
            width: 100%;
            height: 100%;
            :global(img),
            :global(video) {
                object-fit: fill;
                width: 100%;
                max-height: 224px;

                max-width: $max-card-width;
            }
        }
    }
</style>
