import Fetchd from '$elven/network';
import { addTokenToHeaders } from '$elven/tools';
import type { UserChange } from '../types/user';

/** Use with SSR by passing token / or in components by passing empty token.
 *
 * Static methods = not for SSR, use only on components side
 */
export default class NetworkUser {
	private headers: Headers;

	constructor(token: string) {
		const headers = new Headers();
		addTokenToHeaders(token, headers);
		this.headers = headers;
	}

	/** get information about current authorized user */
	public async getMe(): Promise<Response> {
		try {
			const response = await Fetchd.send({ method: 'GET', url: 'users/me', headers: this.headers });
			return response;
		} catch (err) {
			throw err;
		}
	}

	/** change username or password */
	public static async change(body: UserChange) {
		try {
			const resp = await Fetchd.send({ method: 'POST', url: 'users/me/change', body: body });
			return resp;
		} catch (err) {
			return Promise.reject(err);
		}
	}
}
