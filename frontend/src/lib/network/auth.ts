import Fetchd from '$lib/elven/network';
import type { body } from '$lib/types/auth';

/** Not for SSR, use only on components side */
export default class NetworkAuth {
	public static async login(username: string, password: string): Promise<Response> {
		const data: body = {
			username: username,
			password: password,
			type: 'cookie'
		};
		try {
			return await Fetchd.send({ method: 'POST', url: 'auth/login', body: data });
		} catch (err) {
			return Promise.reject(err);
		}
	}

	public static async logout(): Promise<Response> {
		try {
			const resp = await Fetchd.send({ method: 'POST', url: 'auth/logout' });
			return resp;
		} catch (err) {
			return Promise.reject(err);
		}
	}
}
