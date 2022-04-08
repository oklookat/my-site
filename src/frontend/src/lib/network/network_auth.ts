import Fetchd from '$lib/network'
import type { body } from '$lib/types/auth';
import { StorageAuth } from "$lib/tools/storage"

export default class NetworkAuth {

    private headers: Headers

    constructor(token: string) {
        const headers = new Headers()
        StorageAuth.addTokenToHeaders(headers, token)
        this.headers = headers
    }

    public static async login(username: string, password: string): Promise<Response> {
        const data: body = {
            username: username,
            password: password,
            type: 'direct',
        }
        try {
            return await Fetchd.send({ method: "POST", url: 'auth/login', body: data})
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public async logout() {
        try {
            await Fetchd.send({ method: "POST", url: 'auth/logout', headers: this.headers})
            //AuthStorage.set(true)
            return
        } catch (err) {
            return Promise.reject(err)
        }
    }
}