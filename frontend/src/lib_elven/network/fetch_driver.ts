import { NetworkError } from '$elven/network/errors';

export type Fetchable = typeof fetch

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
	driver?: Fetchable;
};

export default class FetchDriver {
	private baseURL = ''

	constructor(baseURL?: string) {
		if(!baseURL) {
			return
		}
		if(!baseURL.endsWith("/")) {
			baseURL += "/"
		}
		this.baseURL = baseURL;
	}

	public async send(config: RequestConfig): Promise<Response> {
		const url = new URL(`${this.baseURL}${config.url}`);

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

		const driver = config.driver || fetch;

		try {
			const result = await driver(url.toString(), fetchConfig);
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
