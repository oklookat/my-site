<script lang="ts">
	import { onMount } from 'svelte';
	import ContextMenu from '$lib/components/context_menu.svelte';
	import NetworkArticle from '$elven/network/article';
	import { isTouchDevice } from '$elven/tools';
	import type { RAW } from '$elven/types/article';
	import { Editable } from '$elven/tools/article';
	import OverlayMobile from '$lib/components/overlay_mobile.svelte';

	/** itself */
	export let article: RAW;

	/** click that triggered this component */
	export let mouseEvent: MouseEvent;

	/** on deleted */
	export let onDeleted: () => void;

	/** on actions closed */
	export let onDisabled: () => void;

	/** choose component to render: overlay / popup */
	const render: {
		active: boolean;
		type: 'context' | 'overlay';
		component: typeof ContextMenu | typeof OverlayMobile | undefined;
		props: any;
	} = { active: true, type: 'overlay', component: undefined, props: null };

	/** is device with touchscreen? */
	let isTouch = false;

	onMount(() => {
		isTouch = isTouchDevice();
		if (isTouch) {
			render.type = 'overlay';
			render.component = OverlayMobile;
			render.props = {
				title: article.title,
				onClose: () => onDisabled()
			};
			return;
		}
		render.type = 'context';
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
			const articled = new Editable(toEdit);
			articled.onSaved = () => {
				onDeleted();
			};
			articled.is_published = isPublished;
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

{#if render.active}
	<svelte:component this={render.component} {...render.props}>
		<div class="actions {render.type}">
			{#if article.is_published}
				<div on:click={async () => await unpublish()}>to drafts</div>
			{:else}
				<div on:click={async () => await publish()}>publish</div>
			{/if}
			<a href={`/elven/articles/create?id=${article.id}`}>edit</a>
			<div on:click={async () => await deleteArticle()}>delete</div>
		</div>
	</svelte:component>
{/if}

<style lang="scss">
	.actions {
		height: 100%;
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		> * {
			cursor: pointer;
			&:hover {
				background-color: var(--color-hover);
			}
			display: flex;
			align-items: center;
			justify-content: center;
		}
		&.overlay,
		&.context {
			> * {
				width: 100%;
			}
		}
		&.overlay {
			padding-top: 12px;
			gap: 14px;
			// item
			> * {
				height: 64px;
				background-color: var(--color-level-1);
			}
		}
		&.context {
			// item
			> * {
				height: 44px;
			}
		}
	}
</style>
