import NetworkAuth from "$lib/network/network_auth";
import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//
import * as cookie from 'cookie';

/** @type {import('./index').RequestHandler} */
export async function post(event: RequestEvent): Promise<RequestHandlerOutput> {
    const body = await event.request.json()
    const username = body.username
    const password = body.password
    //
    const toClient = {
        status: 200,
        headers: {},
        body: {
            loginErr: null
        }
    }
    try {
        const resp = await NetworkAuth.login(username, password)
        if (!resp.ok) {
            toClient.status = resp.status
            if (resp.status === 401) {
                toClient.body.loginErr = 'Incorrect username or password.'
            } else if (resp.status > 500) {
                toClient.body.loginErr = 'Server error. Try later.'
            } else {
                toClient.body.loginErr = 'Unknown error.'
            }
            return toClient
        }

        // we need to get { token: string }
        // json parsing not work idk why
        const data = await resp.json()
        const token = data.token

        // send token cookie
        toClient.headers['Set-Cookie'] = cookie.serialize('token', token, {
            httpOnly: true,
            secure: true,
            maxAge: 63072 * 500, // 1 year
            sameSite: 'none',
            path: '/',
            domain: '.oklookat.ru'
        })
    } catch (err) {
        toClient.status = 400
        toClient.body.loginErr = "Network error. Check your connection."
    }
    return toClient
}