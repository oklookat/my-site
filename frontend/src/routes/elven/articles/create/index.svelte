<script lang="ts" context="module">
	import type { Load, LoadOutput } from '@sveltejs/kit';

	export const load: Load = async (e) => {
		/** creating/editing this article */
		let article: Editable;

		const stuff = e.stuff;
		const output: LoadOutput = {
			status: 200,
			stuff: stuff,
			props: {
				article: article
			}
		};

		const params = e.url.searchParams;
		const isEditMode = params.has('id');
		if (!isEditMode) {
			stuff.title = 'new article';
			return output;
		}

		try {
			const networkArticle = new NetworkArticle(getTokenFromSession(e), e.fetch);
			const resp = await networkArticle.get(params.get('id')!);
			output.status = resp.status;
			if (resp.ok) {
				article = new Editable(await resp.json());
				if (output.props && article.title) {
					output.props.article = article;
					stuff.title = article.title;
				}
			}
		} catch (err) {}
		return output;
	};
</script>

<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';
	import TextareaResizer from '$elven/tools/textarea_resizer';
	import { generateFileTypeSelector } from '$elven/tools/extension';
	import Toolbar from '$lib/components/toolbar.svelte';
	import ArticleCover from '$elven/components/article_cover.svelte';
	import NetworkArticle from '$elven/network/article';
	import FilesPortable from '$elven/components/files_portable.svelte';
	import type { File } from '$elven/types/file';
	import { Params } from '$elven/tools/params';

	import { getEditorConfig } from '$elven/tools/markdown';
	import { getTokenFromSession } from '$elven/tools';
	import UploadPhoto from '$lib/icons/upload_photo.svelte';
	import { Editable } from '$elven/tools/article';
	import { dateToReadable } from '$elven/tools/dates';
	import type jmarkd from '@oklookat/jmarkd';

	/** creating / editing this article */
	export let article: Editable = new Editable();

	/** title element */
	let articleTitleEL: HTMLTextAreaElement;

	/** title resizer */
	let textareaResizer: TextareaResizer;

	/** text editor element */
	let editorEL: HTMLDivElement;

	/** is choose cover overlay opened? */
	let isChooseCover = false;

	/** is cover exists in article? */
	let isCoverExists = false;

	/** md editor instance */
	let editor: jmarkd;

	onMount(async () => {
		// import editor
		// @ts-ignore
		await import('@oklookat/jmarkd/styles');
		const { default: jmarkdModule } = await import('@oklookat/jmarkd');

		if (article) {
			isCoverExists = !!(article.cover_id && article.cover_extension && article.cover_path);
			articleTitleEL.value = article.title;
		}

		textareaResizer = new TextareaResizer(articleTitleEL, 54);
		const config = getEditorConfig(editorEL, article.content);
		editor = new jmarkdModule(config);
	});

	const filesPortableParams = new Params<File>('file');
	filesPortableParams.setParam(
		'extensions',
		generateFileTypeSelector(['IMAGE', 'VIDEO']).selectedToString()
	);

	function onCoverSelected(file: File) {
		isChooseCover = false;
		article.cover_id = file.id;
		article.cover_path = file.path;
		article.cover_extension = file.extension;
		isCoverExists = true;
	}

	function onCoverRemoved() {
		isChooseCover = false;
		isCoverExists = false;
		if (!article.cover_id) {
			return;
		}
		article.cover_id = undefined;
	}

	function onTitleChanged() {
		article.title = articleTitleEL.value;
	}

	function onContentChanged() {
		article.content = editor?.save();
	}

	let lastSavedInterval: NodeJS.Timer;
	let lastSavedPretty = 'not saved';
	article.onSaved = () => {
		if (lastSavedInterval) {
			clearInterval(lastSavedInterval);
		}
		lastSavedInterval = setInterval(() => {
			if (!article.lastSaved) {
				clearInterval(lastSavedInterval);
				return;
			}
			lastSavedPretty = dateToReadable(article.lastSaved) + ' ago';
		}, 1000);
	};

	onDestroy(() => {
		if (lastSavedInterval) {
			clearInterval(lastSavedInterval);
		}
		if (!browser) {
			return;
		}
		editor?.destroy();
		textareaResizer?.destroy();
	});
</script>

{#if isChooseCover}
	<FilesPortable
		params={filesPortableParams}
		on:closed={() => (isChooseCover = false)}
		on:selected={(e) => {
			onCoverSelected(e.detail);
		}}
	>
		<div slot="back-title">article</div>
	</FilesPortable>
{/if}

<div class="create">
	<div class="toolbars">
		<Toolbar>
			<div class="last__saved">
				last saved:
				{lastSavedPretty}
			</div>
		</Toolbar>
	</div>

	<div
		class="cover"
		on:click={() => {
			isChooseCover = !isChooseCover;
		}}
	>
		{#if isCoverExists}
			<div class="remove" on:click|stopPropagation|preventDefault={() => onCoverRemoved()}>X</div>
			<div class="itself">
				<ArticleCover bind:article />
			</div>
		{:else}
			<div class="upload item">
				<UploadPhoto />
			</div>
		{/if}
	</div>

	<div class="editable">
		<textarea
			class="title"
			placeholder="Sample text"
			rows="1"
			maxlength="124"
			bind:this={articleTitleEL}
			on:input={() => onTitleChanged()}
		/>
		<div class="editor" bind:this={editorEL} on:input={() => onContentChanged()} />
	</div>
</div>

<style lang="scss">
	.create {
		max-width: $readable-max-width;
		margin: auto;
		display: flex;
		flex-direction: column;
		gap: 18px;

		.toolbars,
		.editable {
			width: 100%;

			display: flex;
			flex-direction: column;
			gap: 12px;
		}

		.editable {
			width: 100%;
		}

		.cover {
			min-height: 128px;
			background-color: var(--color-level-1);
			position: relative;
			display: flex;
			justify-content: center;

			.itself, .remove {
				cursor: pointer;
			}

			.remove {
				background-color: var(--color-level-1);
				position: absolute;
				width: 48px;
				height: 48px;
				top: 0;
				right: 0;
				z-index: 5;
				display: flex;
				align-items: center;
				justify-content: center;
				align-self: flex-end;

				&:hover {
					background-color: var(--color-hover);
				}
			}

			.upload {
				align-self: center;
				justify-self: center;
				min-height: inherit;
				width: 100%;
				:global(svg) {
					width: 40px;
					height: 40px;

					opacity: 0.5;
					fill: var(--color-text);
				}
			}
		}
	}

	.title,
	.editor {
		box-sizing: border-box;
		border: var(--color-border) 1px solid;
		margin: auto;
		width: 100%;
	}

	.title {
		background-color: var(--color-level-1);
		color: var(--color-text);
		font-size: 1.6rem;
		font-weight: bold;
		min-height: 54px;
		border-radius: 8px;
		padding: 12px;
	}

	.editor {
		min-height: 244px;
		display: flex;
		justify-content: center;
	}
</style>
