import * as path from 'path';
import adapter from '@sveltejs/adapter-node';
import preprocess from 'svelte-preprocess';

const openBrowser = false;
const expose = true;
const port = 3001

// disable cert check when dev.
if (process.env.NODE_ENV == 'development') {
	process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = 0;
}

const viteBase = {
	optimizeDeps: { exclude: [] },
	resolve: {
		alias: {
			$elven: path.resolve('./src/lib_elven'),
			$oklookat: path.resolve('./src/lib_oklookat')
		}
	},
	server: {
		// prevent auto browser opening, because node throws error in container
		host: true,
		open: openBrowser,
		host: expose,
		port: port,
		strictPort: true,
		hmr: {
			// vite HMR WebSocket (SSL) -> nginx -> site.
			protocol: 'wss',
			clientPort: 443
		}
	}
};

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://github.com/sveltejs/svelte-preprocess
	// for more information about preprocessors
	preprocess: preprocess({
		scss: {
			prependData: `@import './src/lib/assets/vars.scss';`
		}
	}),

	kit: {
		vite: viteBase,
		adapter: adapter({ out: './build' }),

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
		csp: {
			mode: 'auto',
			directives: {
				'default-src': [
					'self',
					'data:',
					'*.oklookat.ru',
					'*.google.com',
					'*.youtube.com',
					'*.vk.com',
					'*.github.com',
					'*.facebook.com',
					'*.yandex.ru',
					'*.spotify.com',
					'*.deezer.com',
					'*.twitter.com',
					'*.instagram.com',
					'*.reddit.com'
				]
			}
		}
	}
};

export default config;
