<script lang="ts">
	import { getUploadsWith } from '$lib_elven/tools';
	import Extension, { type FileTypeSelector } from '$lib_elven/tools/extension';
	import type { Article } from '$lib_elven/types/articles';

	export let article: Article;
	$: onArticle(article);

	let isCoverExists = false;
	let extensionSelector: FileTypeSelector;
	let fullPath: string;
	function onArticle(val: Article) {
		isCoverExists = false
		if (!val) {
			return;
		}
		if(!article.cover_id || !article.cover_extension) {
			return
		}
		extensionSelector = Extension.getSelector(article.cover_extension);
		if(!article.cover_path) {
			return
		}
		fullPath = getUploadsWith(article.cover_path).toString();
		isCoverExists = true
	}
</script>

{#if isCoverExists}
	{#if extensionSelector.selected === 'IMAGE'}
		<div class="image">
			<img class="image__main" src={fullPath} alt="cover" />
			<div class="image__blurred" style={`background-image: url(${fullPath})`} />
		</div>
	{:else if extensionSelector.selected === 'VIDEO'}
		<div class="video">
			<video autoplay muted loop src={fullPath} />
		</div>
	{/if}
{/if}

<style lang="scss">
	.image,
	.video {
		width: 100%;
		height: 224px;

		display: flex;
		justify-content: center;
		:global(img),
		:global(video) {
			object-fit: fill;
			width: 100%;
			height: 100%;
		}
	}

	.image,
	.video {
		height: 224px;
		width: 100%;
	}

	.image__main {
		max-width: $desktop-max-card-width;
	}

	.image {
		position: relative;
		overflow: hidden;
		&__main {
			width: 100%;
			position: absolute;
			height: 100%;
			z-index: 7;
			background-size: cover;
		}
		&__blurred {
			width: 100%;
			filter: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='a' x='0' y='0' width='1' height='1'%3E%3CfeGaussianBlur stdDeviation='76' primitiveUnits='userSpaceOnUse' result='b'/%3E%3CfeMorphology operator='dilate' /%3E %3CfeMerge%3E%3CfeMergeNode/%3E%3CfeMergeNode in='b'/%3E%3C/feMerge%3E%3C/filter%3E%3C/svg%3E#a");
			background-repeat: no-repeat;
			background-size: cover;
			background-position: center;
		}
	}
</style>
