import { Axios } from './Axios'
import type { TUser, TUserChange } from '@/types/UserTypes'

export default class UserAdapter {

    public static async getMe(): Promise<TUser> {
        try {
            const response = await Axios.get('/users/me')
            if (response.data) {
                return Promise.resolve(response.data as TUser)
            }
        } catch (err) {
            return Promise.reject(err)
        }
    }


    public static async change(data: TUserChange) {
        try {
            await Axios.post('/users/me/change', data)
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err)
        }
    }

}