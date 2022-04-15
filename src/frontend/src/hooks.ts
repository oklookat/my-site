import type { GetSession, Handle } from '@sveltejs/kit'
//
import NetworkUser from '$lib_elven/network/network_user';
import type { User } from '$lib_elven/types/user';
import Utils from '$lib_elven/tools';
import Validator from '$lib_elven/validators';

export const handle: Handle = async ({ event, resolve }) => {
    let isExists = false
    let isAdmin = false
    let username = ''
    let token = ''
    event.locals.user = {
        isExists: isExists,
        isAdmin: isAdmin,
        username: username,
        token: token
    }

    // get user auth token
    token = Utils.getTokenFromRequestHeaders(event.request.headers)
    let isErr = false

    if (token) {
        const networkUser = new NetworkUser(token)
        let user: User
        try {
            const resp = await networkUser.getMe()
            if (resp.status === 200) {
                user = await resp.json()
            }
        } catch (err) {
            isErr = true
        }
        if (user && user.is_admin && user.username) {
            isExists = true
            isAdmin = user.is_admin
            username = user.username
        }
    }

    event.locals.user.isExists = isExists
    event.locals.user.isAdmin = isAdmin
    event.locals.user.username = username
    event.locals.user.token = token

    const isElvenPage = Validator.isAdminPanelPage(event.url)
    const isElvenLoginPage = Validator.isAdminPanelLoginPage(event.url)
    let redirectURL = `https://${event.url.hostname}`
    if (isElvenPage) {
        if (isErr) {
            const resp = Response.redirect(redirectURL, 302)
            return resp
        }
        if (!isElvenLoginPage && !isAdmin) {
            redirectURL = redirectURL + `/elven/login`
            const resp = Response.redirect(redirectURL, 302)
            return resp
        }
    }

    const response = await resolve(event)
    return response;
}


export const getSession: GetSession = (event) => {
    return {
        user: {
            isExists: event.locals.user.isExists || false,
            isAdmin: event.locals.user.isAdmin || false,
            username: event.locals.user.username || ''
        }
    }
}