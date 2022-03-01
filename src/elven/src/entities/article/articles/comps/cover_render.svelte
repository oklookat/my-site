<script lang="ts">
    import Extension, { type FileType } from "../../../../tools/extension";
    import { PathTools } from "../../../../tools/paths";

    import type { Article } from "../types";

    export let article: Article;
    $: onArticle(article);

    let coverExists = false;
    let extensionType: FileType;
    let fullPath: string;
    function onArticle(val: Article) {
        if (!val) {
            return;
        }
        coverExists = !!(article.cover_id && article.cover_path && article.cover_extension);
        if (!coverExists) {
            return;
        }
        extensionType = Extension.getType(article.cover_extension);
        fullPath = PathTools.getUploadsWith(article.cover_path).toString();
    }
</script>

<div class="cover">
    {#if coverExists}
        {#if extensionType === "image"}
            <div class="cover__image">
                <img
                    decoding="async"
                    loading="lazy"
                    alt="article cover"
                    src={fullPath}
                />
            </div>
        {:else}
            <div class="cover__video">
                <video autoplay muted src={fullPath} />
            </div>
        {/if}
    {/if}
</div>

<style lang="scss">
    .cover{
        width: 100%;
        height: 100%;
        &__image, &__video {
            width: 100%;
            height: 100%;
        }
    }
</style>