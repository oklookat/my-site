import { Axios } from '@/common/adapters/Axios'
import type { IUser } from '@/types/user'

export default class UserAdapter {

    public static getMe() {
        return Axios.get('/users/me')
            .then(response => {
                if (response.data) {
                    return Promise.resolve(response.data as IUser)
                } else {
                    return Promise.reject('No data.')
                }
            })
            .catch(error => {
                return Promise.reject(error)
            })
    }
    
}