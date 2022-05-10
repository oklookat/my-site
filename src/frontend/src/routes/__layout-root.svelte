<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';
	import { locale, loadTranslations } from '$lib/locale';

	export const load: Load = async (event) => {
		const { pathname } = event.url;
		const defaultLocale = 'en';
		const initLocale = locale.get() || defaultLocale;
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
</script>

<slot />
