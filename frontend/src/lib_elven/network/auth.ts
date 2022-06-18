import Fetchd from '$elven/network';
import type { body } from '$elven/types/auth';

/** Not for SSR, use only on components side */
export default class NetworkAuth {
	public static async login(username: string, password: string): Promise<Response> {
		const data: body = {
			username: username,
			password: password,
			type: 'cookie'
		};
		return await Fetchd.send({ method: 'POST', url: 'auth/login', body: data });
	}

	public static async logout(): Promise<Response> {
		return await Fetchd.send({ method: 'POST', url: 'auth/logout' });
	}
}
