<script lang="ts">
	import { getUploadsWith } from '$elven/tools';
	import Extension, { type FileTypeSelector } from '$elven/tools/extension';
	import type { RAW } from '$elven/types/article';

	export let article: RAW;
	$: onArticle(article);

	let isCoverExists = false;
	let extensionSelector: FileTypeSelector;
	let fullPath: string;
	function onArticle(val: RAW) {
		isCoverExists = false;
		if (!val) {
			return;
		}
		if (!article.cover_id || !article.cover_extension) {
			return;
		}
		extensionSelector = Extension.getSelector(article.cover_extension);
		if (!article.cover_path) {
			return;
		}
		fullPath = getUploadsWith(article.cover_path).toString();
		isCoverExists = true;
	}
</script>

{#if isCoverExists}
	{#if extensionSelector.selected === 'IMAGE'}
		<div class="image">
			<div class="blurred" style={`background-image: url(${fullPath})`} />
			<img class="main" src={fullPath} alt="cover" />
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
		max-height: 100%;
		max-width: 100%;
		display: flex;
		justify-content: center;
	}
	.image img,
	.video video {
		width: 100%;
		height: 100%;
		max-height: 480px;
		max-width: 690px;
		object-fit: fill;
	}

	.image {
		position: relative;
		overflow: hidden;
		.main {
			position: relative;
			background-size: cover;
		}
		.blurred {
			width: 100%;
			height: 100%;
			position: absolute;
			filter: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='a' x='0' y='0' width='1' height='1'%3E%3CfeGaussianBlur stdDeviation='76' primitiveUnits='userSpaceOnUse' result='b'/%3E%3CfeMorphology operator='dilate' /%3E %3CfeMerge%3E%3CfeMergeNode/%3E%3CfeMergeNode in='b'/%3E%3C/feMerge%3E%3C/filter%3E%3C/svg%3E#a");
			background-repeat: no-repeat;
			background-size: cover;
			background-position: center;
		}
	}
</style>
