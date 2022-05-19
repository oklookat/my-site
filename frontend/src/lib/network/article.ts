import Fetchd from '$lib/elven/network';
import { addTokenToHeaders } from '$lib/tools';
import type { Article, Params } from '$lib/types/articles';
import type { Operation } from '$lib/tools/jsonpatch';
import type { ArticleRAW } from '$lib/tools/article';

/** Use with SSR by passing token / or in components by passing empty token.
 *
 * Static methods = not for SSR, use only on components side
 */
export default class NetworkArticle {
	private static prefix = 'article/articles';
	private headers: Headers;
	private driver: typeof fetch

	constructor(token: string, driver?: typeof fetch) {
		if (driver) {
			this.driver = driver
		}
		const headers = new Headers();
		addTokenToHeaders(token, headers);
		this.headers = headers;
	}

	public async getAll(params: Params): Promise<Response> {
		const response = await Fetchd.send({
			method: 'GET',
			url: NetworkArticle.prefix,
			params: params,
			headers: this.headers,
			customDriver: this.driver
		});
		return response;
	}

	public async get(id: string): Promise<Response> {
		const response = await Fetchd.send({
			method: 'GET',
			url: `${NetworkArticle.prefix}/${id}`,
			headers: this.headers,
			customDriver: this.driver
		});
		return response
	}

	public static async delete(id: string): Promise<Response> {
		const resp = await Fetchd.send({
			method: 'DELETE',
			url: `${this.prefix}/${id}`
		});
		return resp;
	}

	public static async create(article: ArticleRAW): Promise<Response> {
		const response = await Fetchd.send({
			method: 'POST',
			url: `${this.prefix}`,
			body: article
		});
		return response;
	}

	public static async update(id: string, op: Operation[]): Promise<Response> {
		const response = await Fetchd.send({
			method: 'PATCH',
			url: `${this.prefix}/${id}`,
			body: op
		});
		return response;
	}
}
