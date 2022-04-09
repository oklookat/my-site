import { writable } from 'svelte/store';
//
import * as cookie from 'cookie';

/** app state */
export class StorageGlobal {

    /** is user now in 404 page? */
    public static isNotFoundPage = writable(false);

}

/** auth storage */
export class StorageAuth {

    /** get token from request headers (cookie) */
    public static getToken(headers: Headers): string | null {
        if (!headers || !(headers instanceof Headers)) {
            return null
        }
        if (!headers.has('cookie')) {
            return null
        }
        const cookiesStr = headers.get('cookie')
        const cookiesJson = cookie.parse(cookiesStr)
        if (!cookiesJson || !cookiesJson.token) {
            return null
        }
        return cookiesJson.token
    }

    public static addTokenToHeaders(headers: Headers, token: string) {
        if (!(headers instanceof Headers) || !token) {
            return
        }
        headers.append('Authorization', `Elven ${token}`)
    }

    public static isAdminPanelPage(url: URL): boolean {
        if (!(url instanceof URL)) {
            return false
        }
        const pathname = url.pathname
        return pathname.startsWith("/elven")
    }

    public static isAdminPanelLoginPage(url: URL): boolean {
        if (!(url instanceof URL)) {
            return false
        }
        const pathname = url.pathname
        return pathname.startsWith("/elven/login")
    }
}