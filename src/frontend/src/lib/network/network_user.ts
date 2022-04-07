import Fetchd from '$lib/network'
import { StorageAuth } from '$lib/tools/storage'
import type { User, UserChange } from '../types/user'

export default class NetworkUser {

    private headers: Headers

    constructor(token: string) {
        const headers = new Headers()
        StorageAuth.addTokenToHeaders(headers, token)
        this.headers = headers
    }

    /** get information about current authorized user */
    public async getMe(): Promise<User> {
        try {
            const response = await Fetchd.send({ method: "GET", url: 'users/me', headers: this.headers })
            const jsond = await response.json()
            return jsond as User
        } catch (err) {
            return Promise.reject(err)
        }
    }

    /** change username or password */
    public async change(body: UserChange) {
        try {
            await Fetchd.send({ method: "POST", url: 'users/me/change', body: body, headers: this.headers })
            return
        } catch (err) {
            return Promise.reject(err)
        }
    }

}