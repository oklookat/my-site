<script lang="ts">
	import { onMount } from 'svelte';
	import ContextMenu from '$lib/components/context_menu.svelte';
	import NetworkArticle from '$elven/network/article';
	import { isTouchDevice } from '$elven/tools';
	import type { RAW } from '$elven/types/article';
	import { Editable } from '$elven/tools/article';
import OverlayMobile from '$lib/components/overlay_mobile.svelte';

	export let article: RAW;
	export let mouseEvent: MouseEvent;
	export let onDeleted: () => void;
	export let onDisabled: () => void;

	/** is device with touchscreen? */
	let isTouch = false;

	/** component to render */
	let render: {
		isOverlay: boolean;
		component: typeof ContextMenu | typeof OverlayMobile;
		props: any;
	} = { isOverlay: true, component: OverlayMobile, props: undefined };

	onMount(() => {
		isTouch = isTouchDevice();
		if (isTouch) {
			render.component = OverlayMobile;
			render.props = {
				onClose: () => onDisabled()
			};
			return;
		}
		render.isOverlay = false;
		render.component = ContextMenu;
		render.props = {
			mouseEvent: mouseEvent,
			onDisabled: () => onDisabled()
		};
	});

	async function publishUnpublish(isPublished: boolean) {
		// @ts-ignore
		const toEdit: Article = {
			id: article.id,
			is_published: isPublished
		};
		try {
			// TODO: you know
			const articled = new Editable(toEdit);
			articled.is_published = isPublished;
			// const resp = await NetworkArticle.update(toEdit);
			// if (resp.ok) {
			// 	onDeleted();
			// }
		} catch (err) {
			return Promise.reject(err);
		}
	}

	/** publish article */
	async function publish() {
		await publishUnpublish(true);
	}

	/** unpublish article */
	async function unpublish() {
		await publishUnpublish(false);
	}

	/** delete article */
	async function deleteArticle() {
		if (!window.$confirm || !article.id) {
			return;
		}
		onDisabled();
		const isDelete = await window.$confirm('Delete article');
		if (!isDelete) {
			return;
		}
		try {
			await NetworkArticle.delete(article.id);
			onDeleted();
		} catch (err) {}
	}
</script>

<svelte:component this={render.component} {...render.props}>
	{#if article.is_published}
		<div on:click={async () => await unpublish()}>to drafts</div>
	{:else}
		<div on:click={async () => await publish()}>publish</div>
	{/if}
	<a href={`/elven/articles/create?id=${article.id}`}>edit</a>
	<div on:click={async () => await deleteArticle()}>delete</div>
</svelte:component>
