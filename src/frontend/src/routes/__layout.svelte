<script context="module" lang="ts">
	import { browser } from '$app/env';
	import type { Load } from '@sveltejs/kit';
	import { getLocaleFromNavigator, init } from 'svelte-i18n';
	import { register } from 'svelte-i18n';

    // https://github.com/kaisermann/svelte-i18n/issues/166#issuecomment-1001009977
	register('en', () => import('../lib/locales/en.json'));
	register('ru', () => import('../lib/locales/ru.json'));

	if (browser) {
		init({
			fallbackLocale: 'en',
			initialLocale: getLocaleFromNavigator()
		});
	}

	export const load: Load = (event) => {
		if (!browser) {
			init({
				fallbackLocale: 'en',
				initialLocale: 'en'
			});
		}

		return {};
	};
</script>

<slot />
