<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	// ui
	import Overlay from '$lib/components/overlay.svelte';
	import Popup from '$lib/components/popup.svelte';
	// utils
	import Validator from '$lib/validators';
	// article
	import type { Article } from '$lib/types/articles';

	/** file itself */
	export let article: Article;

	/** click on file mouse event */
	export let mouseEvent: MouseEvent;

	/** on file deleted */
	export let onDeleted: () => void;

	/** on actions closed */
	export let onDisabled: () => void;

	/** is device with touchscreen? */
	const isTouchDevice = Validator.isTouchDevice();

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

	async function publishUnpublish(isPublished: boolean): Promise<Response> {
		// @ts-ignore
		const toEdit: Article = {
			id: article.id,
			is_published: isPublished
		};
		try {
			const resp = await fetch('/elven/articles/create', {
				method: 'PATCH',
				body: JSON.stringify(toEdit)
			});
			return Promise.resolve(resp);
		} catch (err) {
			return Promise.reject(err);
		}
	}

	/** publish article */
	async function publish() {
		publishUnpublish(true);
	}

	/** unpublish article */
	async function unpublish() {
		publishUnpublish(false);
	}

	/** edit article */
	function edit() {
		goto(`/elven/articles/create?id=${article.id}`);
	}

	/** delete article */
	async function deleteArticle() {
		const isDelete = await window.$choose.confirm('delete article');
		if (!isDelete) {
			return;
		}
		try {
			await fetch('/elven/articles/create', {
				method: 'DELETE',
				body: article.id
			});
			onDeleted();
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
