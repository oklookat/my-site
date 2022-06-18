//
import Fetchd from '$elven/network';
import type { Fetchable } from '$elven/network/fetch_driver';
import { addTokenToHeaders } from '$elven/tools';
import type { Params } from '$elven/types/file';

/** Use with SSR by passing token / or in components by passing empty token.
 *
 * Static methods = not for SSR, use only on components side
 */
export default class NetworkFile {
	private headers: Headers;

	constructor(token: string) {
		const headers = new Headers();
		addTokenToHeaders(token, headers);
		this.headers = headers;
	}

	/** get files list */
	public async getAll(params: Params, driver?: Fetchable): Promise<Response> {
		return await Fetchd.send({
			method: 'GET',
			url: 'files',
			params: params,
			headers: this.headers,
			driver: driver
		});
	}

	/** upload one file */
	public static async upload(file: File): Promise<Response> {
		if (!(file instanceof File)) {
			throw Error('not a File');
		}
		const form = new FormData();
		form.append('file', file);
		return await Fetchd.send({ method: 'POST', url: 'files', body: form });
	}

	/** delete one file */
	public static async delete(id: string): Promise<Response> {
		return await Fetchd.send({ method: 'DELETE', url: `files/${id}` });
	}
}
