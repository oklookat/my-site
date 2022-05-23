<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';
	import { loadTranslations } from '$lib/locale';

	export const load: Load = async (event) => {
		const { pathname } = event.url;
		const defaultLocale = 'en';
		let browserLocale: string | null = null;
		if (browser) {
			browserLocale = navigator.language;
		}
		const initLocale = browserLocale || defaultLocale;
		await loadTranslations(initLocale, pathname);
		return {};
	};
</script>

<script lang="ts">
	// main style
	import '$lib/assets/base.scss';
	import '$lib/assets/fonts.css';
	import '$lib/assets/global.scss';
	import '$lib/assets/root.css';
	import '$lib/assets/ui.scss';
	import { browser } from '$app/env';
</script>

<slot />
