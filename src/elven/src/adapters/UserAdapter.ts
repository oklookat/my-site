import Duck from './Duck'
import type { TUser, TUserChange } from '@/types/UserTypes'

export default class UserAdapter {

    public static async getMe(): Promise<TUser> {
        try {
            const response = await Duck.GET({ url: '/users/me' })
            return Promise.resolve(response.body as TUser)
        } catch (err) {
            return Promise.reject(err)
        }
    }


    public static async change(body: TUserChange) {
        try {
            await Duck.POST({ url: '/users/me/change', body })
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err)
        }
    }

}