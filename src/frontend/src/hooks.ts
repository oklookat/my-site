import type { GetSession, Handle } from '@sveltejs/kit'
//
import NetworkUser from '$lib/network/network_user';
import { StorageAuth } from '$lib/tools/storage';
import type { User } from '$lib/types/user';

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
    token = StorageAuth.getToken(event.request.headers)
    let isErr = false

    if (token) {
        const networkUser = new NetworkUser(token)
        let user: User
        try {
            user = await networkUser.getMe()
        } catch(err) {
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

    const isElvenPage = StorageAuth.isAdminPanelPage(event.url)
    const isElvenLoginPage = StorageAuth.isAdminPanelLoginPage(event.url)
    if (!isErr && isElvenPage && !isElvenLoginPage && !isAdmin) {
        //const resp = Response.redirect('https://www.youtube.com/watch?v=dQw4w9WgXcQ', 302)
        const resp = Response.redirect('https://oklookat.ru/elven/login', 302)
        return resp
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