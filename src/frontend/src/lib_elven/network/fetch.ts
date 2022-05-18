import { NetworkError } from '$lib_elven/network/errors';

export type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

export interface Stringer {
	toString: () => string;
}
export type RequestConfig = {
	method: RequestMethod;
	url: string;
	body?: BodyInit | Object;
	params?: Record<string | number, Stringer>;
	headers?: Headers;
	customDriver?: typeof fetch;
};

export default class FetchDriver {
	private baseURL: string;

	constructor(baseURL: string) {
		this.baseURL = baseURL;
	}

	public async send(config: RequestConfig): Promise<Response> {
		let endpoint = '';

		// set request url
		if (this.baseURL) {
			endpoint = `${this.baseURL}/${config.url}`;
		} else {
			endpoint = config.url;
		}
		const url = new URL(endpoint);
		if (config.params) {
			// add searchparams to url
			for (const key in config.params) {
				const val = config.params[key];
				const isValInvalid = !val && typeof val !== 'boolean';
				if (isValInvalid) {
					continue;
				}
				// @ts-ignore
				url.searchParams.append(key, val);
			}
		}

		const headers = new Headers();
		headers.append('Accept', 'application/json');
		headers.append('Content-Type', 'application/json');

		// actions depend on request body type
		if (config.body) {
			const isFormData = typeof FormData !== 'undefined' && config.body instanceof FormData;
			if (isFormData) {
				headers.delete('Content-Type');
			} else {
				config.body = JSON.stringify(config.body || {});
			}
		}

		if (config.headers) {
			config.headers.forEach((val, key) => {
				if (!val || key.toUpperCase() == 'CONTENT-TYPE') {
					return;
				}
				headers.append(key, val);
			});
		}

		const fetchConfig: RequestInit = {
			method: config.method,
			credentials: 'include',
			headers: headers,
			// @ts-ignore
			body: config.body
		};

		if (config.method == 'GET') {
			delete fetchConfig.body;
		}

		let okFetch = fetch;
		if (config.customDriver) {
			okFetch = config.customDriver;
		}

		try {
			const result = await okFetch(url.toString(), fetchConfig);
			if (!result.ok) {
				NetworkError.handle(result);
			}
			return result;
		} catch (err) {
			NetworkError.handle(err as any);
			return Promise.reject(err);
		}
	}
}
