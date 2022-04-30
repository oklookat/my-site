<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';
	import { marked } from 'marked';
	import hljs from 'highlight.js';
	import '../../../../lib_elven/assets/highlight.scss';
	import type { Config } from '@oklookat/jmarkd';
	import TextareaResizer from '$lib_elven/tools/textarea_resizer';
	import { generateFileTypeSelector } from '$lib_elven/tools/extension';
	import Toolbar from '$lib_elven/components/toolbar.svelte';
	import type { Article } from '$lib_elven/types/articles';
	import ArticleCover from '$lib_elven/components/article_cover.svelte';
	import NetworkArticle from '$lib_elven/network/network_article';
	import FilesPortable from '$lib_elven/components/files_portable.svelte';
	import type { File } from '$lib_elven/types/files';
	import { setTitleElven } from '$lib_elven/tools';
	import ToolsArticles from '$lib_elven/tools/articles';
	import { Params } from '$lib_elven/tools/params';
	import { dateToReadable } from '$lib_elven/tools/dates';
	import { _ } from 'svelte-i18n'

	/** creating / editing this article */
	export let article: Article;

	/** save all data */
	const save = saver();

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
	function onCoverChanged(val: string | undefined) {
		isCoverExists = !!(article.cover_id && article.cover_path && article.cover_extension);
	}

	/** md editor */
	let jmarkdClass: any;

	/** md editor instance */
	let editor: any;

	onMount(async () => {
		// import markdown editor
		// styles
		// @ts-ignore
		await import('@oklookat/jmarkd/styles');
		// editor class
		const { default: jmarkdModule } = await import('@oklookat/jmarkd');
		jmarkdClass = jmarkdModule;

		// manually add title before creating TextareaResizer, for correct height in start
		textareaResizer = new TextareaResizer(articleTitleEL, 54);
		articleTitleEL.value = article.title;
		initEditor(article.content);
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		editor?.destroy();
		textareaResizer?.destroy();
	});

	/** start text editor */
	function initEditor(data?: string) {
		marked.setOptions({
			renderer: new marked.Renderer(),
			highlight: function (code, lang) {
				const language = hljs.getLanguage(lang) ? lang : 'plaintext';
				return hljs.highlight(code, { language }).value;
			},
			// highlight.js css expects a top-level 'hljs' class.
			langPrefix: 'hljs language-',
			pedantic: false,
			gfm: true,
			breaks: false,
			sanitize: false,
			smartLists: true,
			smartypants: false,
			xhtml: false
		});
		const config: Config = {
			container: editorEL,
			placeholder: $_('elven.routes.articles.create.editorPlaceholder'),
			input: data,
			toolbar: {
				elements: {
					config: {
						preview: {
							parse: (data: string) => {
								return marked.parse(data);
							}
						}
					}
				}
			}
		};
		editor = new jmarkdClass(config);
	}

	let lastSavedPretty = $_('elven.routes.articles.create.notSaved');
	const updateLastSaved = createLastSaver();
	function createLastSaver() {
		let lastSavedTimestamp = 0;
		let lastSavedInterval: NodeJS.Timer;
		return () => {
			lastSavedTimestamp = new Date().getTime();
			lastSavedPretty = dateToReadable(lastSavedTimestamp);

			if (lastSavedInterval) {
				clearInterval(lastSavedInterval);
			}

			lastSavedInterval = setInterval(() => {
				lastSavedPretty = dateToReadable(lastSavedTimestamp);
			}, 1000);
		};
	}

	/** create new article */
	async function createArticle(): Promise<Article> {
		const notValid =
			article.id ||
			!ToolsArticles.validateTitle(article.title) ||
			!ToolsArticles.validateContent(article.content);
		if (notValid) {
			return Promise.reject('not valid article');
		}
		window.$progress?.startBasic();
		try {
			const resp = await NetworkArticle.create(article);
			window.$progress?.finishBasic();
			if (resp.ok) {
				const newArticle = await resp.json();
				article.id = newArticle.id;
				updateLastSaved();
				return newArticle;
			}
			throw resp.statusText;
		} catch (err) {
			throw err;
		} finally {
		}
	}

	/** update existing article */
	async function updateArticle(): Promise<Article> {
		const notValid =
			!article.id ||
			!ToolsArticles.validateTitle(article.title) ||
			!ToolsArticles.validateContent(article.content);
		if (notValid) {
			return Promise.reject("not valid article");
		}
		window.$progress?.startBasic();
		const resp = await NetworkArticle.update(article);
		window.$progress?.finishBasic();
		if (resp.ok) {
			updateLastSaved();
			return await resp.json();
		}
		throw Error(resp.statusText);
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

	const filesPortableParams = new Params<File>('file');
	filesPortableParams.setParam(
		'extensions',
		generateFileTypeSelector(['IMAGE', 'VIDEO']).selectedToString()
	);

	const createArticleTitle = $_('elven.routes.articles.create.createArticle')
</script>

<svelte:head>
	<title>{setTitleElven(`${article.id ? article.title : createArticleTitle}`)}</title>
</svelte:head>

{#if isChooseCover}
	<FilesPortable
		params={filesPortableParams}
		on:closed={() => (isChooseCover = false)}
		on:selected={(e) => {
			onCoverSelected(e.detail);
		}}
	>
		<div slot="back-title">{$_('elven.general.article')}</div>
	</FilesPortable>
{/if}

<div class="create base__container">
	<div class="toolbars">
		<Toolbar>
			<div class="last__saved">
				{$_('elven.routes.articles.create.lastSaved')} {lastSavedPretty}
			</div>
		</Toolbar>
	</div>

	<div
		class="cover pointer with-border"
		on:click={() => {
			isChooseCover = !isChooseCover;
		}}
	>
		{#if isCoverExists}
			<div class="remove" on:click|stopPropagation|preventDefault={() => removeCover()}>X</div>
			<div class="itself">
				<ArticleCover bind:article />
			</div>
		{:else}
			<div class="upload item">
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
			placeholder={$_('elven.routes.articles.create.titlePlaceholder')}
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
		display: flex;
		flex-direction: column;
		gap: 18px;

		max-width: $readable-max-width;
		margin: auto;
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
	}

	.cover {
		position: relative;
		// overflow for clipping child border
		overflow: hidden;
		background-color: var(--color-level-1);

		display: flex;
		flex-direction: column;

		width: 100%;
		height: 224px;
		.remove {
			width: 48px;
			height: 48px;

			z-index: 777;

			top: 0;
			right: 0;

			background-color: inherit;

			display: flex;
			align-items: center;
			justify-content: center;
			align-self: flex-end;

			border-bottom: var(--border);
			border-left: var(--border);

			&:hover {
				background-color: var(--color-hover);
			}
		}
		.upload {
			width: 100%;
			svg {
				width: 40px;
				height: 40px;

				opacity: 0.5;
				fill: var(--color-text);
			}
		}
		.itself {
			position: absolute;
			height: max-content;
			width: 100%;
			display: flex;
			justify-content: center;
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
