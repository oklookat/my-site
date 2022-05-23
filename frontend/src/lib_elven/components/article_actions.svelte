<script lang="ts">
	import { onMount } from 'svelte';
	import Overlay from '$elven/components/overlay.svelte';
	import Popup from '$elven/components/popup.svelte';
	import NetworkArticle from '$elven/network/article';
	import { isTouchDevice } from '$elven/tools';
	import { t } from '$lib/locale';
	import type { RAW } from '$elven/types/article';
	import { Editable } from '$elven/tools/article';

	export let article: RAW;
	export let mouseEvent: MouseEvent;
	export let onDeleted: () => void;
	export let onDisabled: () => void;

	/** is device with touchscreen? */
	let isTouch = false;

	/** component to render */
	let render: {
		isOverlay: boolean;
		component: any;
		props: any;
	} = { isOverlay: true, component: null, props: null };

	onMount(() => {
		isTouch = isTouchDevice();
		if (isTouch) {
			render.component = Overlay;
			render.props = {
				onClose: () => onDisabled()
			};
			return;
		}
		render.isOverlay = false;
		render.component = Popup;
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
		const isDelete = await window.$confirm($t('elven.articles.deleteQuestion'));
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
		<div on:click={async () => await unpublish()}>
			{$t('elven.articles.toDrafts')}
		</div>
	{:else}
		<div on:click={async () => await publish()}>
			{$t('elven.articles.publish')}
		</div>
	{/if}
	<a href={`/elven/articles/create?id=${article.id}`}>{$t('elven.articles.edit')}</a>
	<div on:click={async () => await deleteArticle()}>{$t('elven.articles.delete')}</div>
</svelte:component>
