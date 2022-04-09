<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';
	// editor
	import type { config } from '@oklookat/jmarkd';

	// utils
	import Animation from '$lib/tools/animation';
	import TextareaResizer from '$lib/tools/textarea_resizer';
	import { generateFileTypeSelector } from '$lib/tools/extension';

	// ui
	import Toolbar from '$lib/components/toolbar.svelte';

	// article
	import type { Article } from '$lib/types/articles';
	import ValidatorArticle from '$lib/validators/validator_article';
	import CategoriesSelector from '$lib/components/categories_selector.svelte';
	import type { Category } from '$lib/types/articles/categories';
	import ArticleCover from '$lib/components/article_cover.svelte';
	import NetworkArticle from '$lib/network/network_article';

	// file
	import FilesPortable from '$lib/components/files_portable.svelte';
	import type { File } from '$lib/types/files';

	/** creating / editing this article */
	export let article: Article;

	/** save all data */
	const save = saver();

	/** main container */
	let createContainer: HTMLDivElement;

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

	$: onCoverChanged(article.cover_id);
	function onCoverChanged(val) {
		isCoverExists = !!(article.cover_id && article.cover_path && article.cover_extension);
	}

	/** md editor */
	let jmarkdClass;

	/** md editor instance */
	let editor;

	onMount(async () => {
		// import markdown editor
		const jmarkdModule = await import('@oklookat/jmarkd');
		jmarkdClass = jmarkdModule.default;
		// @ts-ignore
		await import('@oklookat/jmarkd/styles');

		// manually add title before creating TextareaResizer, for correct height in start
		textareaResizer = new TextareaResizer(articleTitleEL, 54);
		articleTitleEL.value = article.title;
		initEditor(article.content);

		// all loaded - set opacity
		// (not display, because it brokes title resizing on init)
		await Animation.fadeIn(createContainer);
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}

		if (editor) {
			editor.destroy();
		}

		if (textareaResizer) {
			textareaResizer.destroy();
		}
	});

	/** start text editor */
	function initEditor(data?: string) {
		// TODO: add sanitizer
		const config: config = {
			container: editorEL,
			placeholder: `It's a long story...`,
			input: data
		};
		editor = new jmarkdClass(config);
	}

	/** create new article */
	async function createArticle() {
		const notValid =
			article.id ||
			!ValidatorArticle.title(article.title) ||
			!ValidatorArticle.content(article.content);
		if (notValid) {
			return;
		}
		try {
			const newArticle = await NetworkArticle.create(article);
			article.id = newArticle.id;
		} catch (err) {
			return Promise.reject(err);
		}
		return Promise.resolve();
	}

	/** update existing article */
	async function updateArticle() {
		const notValid =
			!article.id ||
			!ValidatorArticle.title(article.title) ||
			!ValidatorArticle.content(article.content);
		if (notValid) {
			return;
		}
		const updated = await NetworkArticle.update(article);
		return updated;
	}

	/** create save func */
	function saver() {
		let throttle: NodeJS.Timeout;

		// save logic
		const saver = async () => {
			const outputData = editor.save();
			article.content = outputData;
			// if saved before (update)
			if (article.id) {
				return await updateArticle();
			}
			// if not saved before (new article)
			return await createArticle();
		};

		// save data
		return (): Promise<null> => {
			if (throttle) {
				clearTimeout(throttle);
			}
			return new Promise((resolve, reject) => {
				throttle = setTimeout(async () => {
					try {
						await saver();
						resolve(null);
					} catch (err) {
						reject(err);
					}
				}, 1000);
			});
		};
	}

	function onCategoryChanged(newCat: Category | null) {
		// category not changed?
		const newCatNotEmpty = newCat && newCat.id;
		const notChanged =
			(!article.category_id && !newCat) || (newCatNotEmpty && article.category_id === newCat.id);
		if (notChanged) {
			return;
		}
		const oldCatID = article.category_id;
		if (newCatNotEmpty) {
			article.category_id = newCat.id;
		} else {
			// no category
			article.category_id = null;
		}
		save().catch(() => {
			// revert changes
			article.category_id = oldCatID;
		});
	}

	function onCoverSelected(file: File) {
		isChooseCover = false;
		article.cover_id = file.id;
		article.cover_extension = file.extension;
		article.cover_path = file.path;
		save();
	}

	function removeCover() {
		isChooseCover = false;
		if (!article.cover_id) {
			return;
		}
		article.cover_id = undefined;
		save();
	}
</script>

<svelte:head>
	<title>{`elven: ${article.id ? article.title : 'create article'}`}</title>
</svelte:head>

{#if isChooseCover}
	<FilesPortable
		params={{
			extensions: generateFileTypeSelector(['IMAGE', 'VIDEO']).selectedToString()
		}}
		on:closed={() => (isChooseCover = false)}
		on:selected={(e) => {
			onCoverSelected(e.detail);
		}}
	/>
{/if}

<div class="create base__container" bind:this={createContainer}>
	<div class="toolbars">
		<Toolbar>
			<CategoriesSelector
				bind:selectedID={article.category_id}
				on:changed={(e) => onCategoryChanged(e.detail)}
			/>
			{#if isCoverExists}
				<div class="remove-cover button" on:click={() => removeCover()}>remove cover</div>
			{/if}
		</Toolbar>
	</div>

	<div
		class="cover pointer"
		on:click={() => {
			isChooseCover = !isChooseCover;
		}}
	>
		{#if isCoverExists}
			<div class="cover__itself">
				<ArticleCover bind:article />
			</div>
		{:else}
			<div class="cover__upload item with-border">
				<svg
					version="1.1"
					xmlns="http://www.w3.org/2000/svg"
					xmlns:xlink="http://www.w3.org/1999/xlink"
					x="0px"
					y="0px"
					viewBox="0 0 230 230"
					style="enable-background:new 0 0 230 230;"
					xml:space="preserve"
				>
					<path
						d="M132.651,140.748H97.349v-35.301h35.302V140.748z M59.32,52.496H230v141.203H0V52.496h17.571V36.301H59.32V52.496z
      M166.313,81.975h45.491V67.781h-45.491V81.975z M65.87,123.096c0,27.136,21.996,49.131,49.13,49.131s49.13-21.995,49.13-49.131
     c0-27.131-21.996-49.129-49.13-49.129S65.87,95.965,65.87,123.096z"
					/>
				</svg>
			</div>
		{/if}
	</div>

	<div class="editable">
		<textarea
			class="title"
			placeholder="Actually..."
			rows="1"
			maxlength="124"
			bind:value={article.title}
			bind:this={articleTitleEL}
			on:input={() => save()}
		/>
		<div class="editor" bind:this={editorEL} on:input={() => save()} />
	</div>
</div>

<style lang="scss">
	.create {
		// after data loaded - opacity = 1
		opacity: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 18px;
		.toolbars,
		.editable {
			width: 100%;
			display: flex;
			flex-direction: column;
			gap: 12px;
		}
		.cover,
		.editable {
			width: 100%;
			max-width: 744px;
		}
	}

	.remove-cover {
		width: 104px;
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.cover {
		background-color: var(--color-level-1);
		display: flex;
		justify-content: center;
		justify-items: center;
		&__upload {
			width: 100%;
			height: 84px;
			svg {
				width: 40px;
				height: 40px;
				opacity: 0.5;
				@media screen and(prefers-color-scheme: light) {
					fill: black;
				}
				@media screen and(prefers-color-scheme: dark) {
					fill: white;
				}
			}
		}
		&__itself {
			width: 100%;
			height: max-content;
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
		background-color: white;
		color: black;
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
