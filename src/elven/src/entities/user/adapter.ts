import Duck from '@/duck'
import type { User, UserChange } from './types'

export default class UserAdapter {

    public static async getMe(): Promise<User> {
        try {
            const response = await Duck.GET({ url: '/users/me' })
            return Promise.resolve(response.body as User)
        } catch (err) {
            return Promise.reject(err)
        }
    }


    public static async change(body: UserChange) {
        try {
            await Duck.POST({ url: '/users/me/change', body })
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err)
        }
    }

}