<script lang="ts">
    import { onMount } from "svelte";
    import { push } from "svelte-spa-router";
    // ui
    import Overlay from "@/components/overlay.svelte";
    import Popup from "@/components/popup.svelte";
    // utils
    import Utils from "@/tools";
    // article
    import NetworkArticle from "@/network/network_article";
    import type { Article } from "@/types/articles";


    /** file itself */
    export let article: Article;

    /** click on file mouse event */
    export let mouseEvent: MouseEvent;

    /** on file deleted */
    export let onDeleted: () => void;

    /** on actions closed */
    export let onDisabled: () => void;

    /** is device with touchscreen? */
    const isTouchDevice = Utils.isTouchDevice();

    /** component to render */
    let render: {
        isOverlay: boolean;
        component: any;
        props: any;
    } = { isOverlay: true, component: null, props: null };

    onMount(() => {
        if (isTouchDevice) {
            render.component = Overlay;
            render.props = {
                onClose: () => onDisabled(),
            };
            return;
        }
        render.isOverlay = false;
        render.component = Popup;
        render.props = {
            mouseEvent: mouseEvent,
            onDisabled: () => onDisabled(),
        };
    });

    /** unpublish article */
    async function unpublish() {
        try {
            await NetworkArticle.unpublish(article.id);
            onDeleted()
        } catch (err) {}
    }

    /** publish article */
    async function publish() {
        try {
            await NetworkArticle.publish(article.id);
            onDeleted()
        } catch (err) {}
    }

    /** edit article */
    function edit() {
        push(`/articles/create/${article.id}`);
    }

    /** delete article */
    async function deleteArticle() {
        const isDelete = await window.$choose.confirm("delete article");
        if (!isDelete) {
            return;
        }
        try {
            await NetworkArticle.delete(article.id);
            onDeleted()
        } catch (err) {}
    }
</script>

<svelte:component this={render.component} {...render.props}>
    <div class="base__items {render.isOverlay ? 'extended' : ''}">
        {#if article.is_published}
            <div on:click={() => unpublish()}>unpublish</div>
        {:else}
            <div on:click={() => publish()}>publish</div>
        {/if}
        <div on:click={() => edit()}>edit</div>
        <div on:click={() => deleteArticle()}>delete</div>
    </div>
</svelte:component>
