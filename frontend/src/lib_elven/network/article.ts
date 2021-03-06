import Fetchd from '$elven/network';
import { addTokenToHeaders } from '$elven/tools';
import type { Params } from '$elven/types/article';
import type { Operation } from '$elven/tools/jsonpatch';
import type { RAW } from '$elven/types/article';
import type { Fetchable } from '$elven/network/fetch_driver';

/** Use with SSR by passing token / or in components by passing empty token.
 *
 * Static methods = not for SSR, use only on components side
 */
export default class NetworkArticle {
	private static prefix = 'article/articles';
	private headers: Headers;
	private driver: Fetchable;

	constructor(token: string, driver?: Fetchable) {
		if (driver) {
			this.driver = driver;
		}
		const headers = new Headers();
		addTokenToHeaders(token, headers);
		this.headers = headers;
	}

	public async getAll(params: Params): Promise<Response> {
		return await Fetchd.send({
			method: 'GET',
			url: NetworkArticle.prefix,
			params: params,
			headers: this.headers,
			driver: this.driver
		});
	}

	public async get(id: string): Promise<Response> {
		return await Fetchd.send({
			method: 'GET',
			url: `${NetworkArticle.prefix}/${id}`,
			headers: this.headers,
			driver: this.driver
		});
	}

	public static async delete(id: string): Promise<Response> {
		return await Fetchd.send({
			method: 'DELETE',
			url: `${this.prefix}/${id}`
		});
	}

	public static async create(article: RAW): Promise<Response> {
		return await Fetchd.send({
			method: 'POST',
			url: `${this.prefix}`,
			body: article
		});
	}

	public static async update(id: string, op: Operation[]): Promise<Response> {
		return await Fetchd.send({
			method: 'PATCH',
			url: `${this.prefix}/${id}`,
			body: op
		});
	}
}
