import type { GetSession, Handle } from '@sveltejs/kit';
//
import NetworkUser from '$elven/network/user';
import type { User } from '$elven/types/user';
import { getTokenFromRequestHeaders } from '$elven/tools';

export const handle: Handle = async ({ event, resolve }) => {
	const resolver = async (): Promise<Response> => {
		const resp = await resolve(event);
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Feature-Policy
		resp.headers.append(
			'Permissions-Policy',
			'microphone=(), geolocation=(), camera=(), payment=(), usb=()'
		);
		return resp;
	};

	let isError = false;
	let isExists = false;
	let isAdmin = false;
	let username = '';
	let token: string | null = '';
	event.locals.user = {
		isError: isError,
		isExists: isExists,
		isAdmin: isAdmin,
		username: username,
		token: token
	};

	// get user auth token
	token = getTokenFromRequestHeaders(event.request.headers);

	if (!token) {
		return await resolver();
	}

	const networkUser = new NetworkUser(token);
	let user: User | undefined = undefined;
	try {
		const resp = await networkUser.getMe();
		if (resp.ok) {
			user = await resp.json();
		} else {
			isError = true;
		}
	} catch (err) {
		isError = true;
	}

	if (!isError && user && user.is_admin && user.username) {
		isExists = true;
		isAdmin = user.is_admin;
		username = user.username;
	}

	event.locals.user.isError = isError;
	event.locals.user.isExists = isExists;
	event.locals.user.isAdmin = isAdmin;
	event.locals.user.username = username;
	event.locals.user.token = token;

	return await resolver();
};

export const getSession: GetSession = (event) => {
	return {
		user: {
			isError: event.locals.user.isError || false,
			isExists: event.locals.user.isExists || false,
			isAdmin: event.locals.user.isAdmin || false,
			username: event.locals.user.username || '',
			token: event.locals.user.token || ''
		}
	};
};
