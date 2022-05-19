<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Overlay from '$lib_elven/components/overlay.svelte';
	import Popup from '$lib_elven/components/popup.svelte';
	import type { Article } from '$lib_elven/types/articles';
	import NetworkArticle from '$lib_elven/network/network_article';
	import { isTouchDevice } from '$lib/tools';
	import { t } from '$lib/locale';

	export let article: Article;
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
			const resp = await NetworkArticle.update(toEdit);
			if (resp.ok) {
				onDeleted();
			}
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

	/** edit article */
	async function edit() {
		await goto(`/elven/articles/create?id=${article.id}`);
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
	<div on:click={async () => await edit()}>{$t('elven.articles.edit')}</div>
	<div on:click={async () => await deleteArticle()}>{$t('elven.articles.delete')}</div>
</svelte:component>
