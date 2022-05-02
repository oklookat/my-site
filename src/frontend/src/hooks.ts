import type { GetSession, Handle } from '@sveltejs/kit';
//
import NetworkUser from '$lib_elven/network/network_user';
import type { User } from '$lib_elven/types/user';
import { getTokenFromRequestHeaders } from '$lib/tools';

export const handle: Handle = async ({ event, resolve }) => {
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
		const response = await resolve(event);
		return response;
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

	const response = await resolve(event);
	return response;
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
