import Duck from '@/network'
import type { User, UserChange } from '../types/user'

export default class UserNetwork {

    /** get information about current authorized user */
    public static async getMe(): Promise<User> {
        try {
            const response = await Duck.GET({ url: 'users/me' })
            return Promise.resolve(response.body as User)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    /** change username or password */
    public static async change(body: UserChange) {
        try {
            await Duck.POST({ url: 'users/me/change', body })
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err)
        }
    }

}