import * as path from 'path';
import adapter from '@sveltejs/adapter-node';
import preprocess from 'svelte-preprocess';

const openBrowser = false;
const expose = false;

const viteBase = {
	optimizeDeps: { exclude: [] },
	resolve: {
		alias: {
			$lib_elven: path.resolve('./src/lib_elven'),
			$lib_oklookat: path.resolve('./src/lib_oklookat')
		}
		// include .d.ts
		//extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.d.ts']
	},
	server: {
		// prevent auto browser opening, because node throws error in container
		host: true,
		open: openBrowser,
		host: expose,
		port: 3000,
		strictPort: true,
		hmr: {
			// vite HMR WebSocket (SSL) -> nginx -> site.
			protocol: 'wss',
			clientPort: 443
		}
	}
};

if (process.env.NODE_ENV == 'development') {
	process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = 0;
}

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://github.com/sveltejs/svelte-preprocess
	// for more information about preprocessors
	preprocess: preprocess({
		scss: {
			prependData: `@import './src/lib_elven/assets/utils.scss';`
		}
	}),

	kit: {
		adapter: adapter({ out: './build' }),
		vite: viteBase
	}
};

export default config;
