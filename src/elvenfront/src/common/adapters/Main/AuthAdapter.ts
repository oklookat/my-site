import Axios from '@/common/adapters/Axios.js'
import Store from "@/store/index"

class AuthAdapter {

    public static async login(username: string, password: string) {
        const data = {
            username: username,
            password: password,
            type: 'admin',
        }
        return await Axios.post('auth/login', data)
            .then(async response => {
                const token = response.data.token
                if (token) {
                    await Store.dispatch('setToken', token)
                    return Promise.resolve()
                }
                return Promise.reject('Ошибка при получении токена.')
            })
            .catch(error => {
                if (!error.response) {
                    return Promise.reject('Произошла ошибка. Попробуйте позже.')
                }
                const readableErr = error.response.data.error
                if(readableErr){
                    return Promise.reject(readableErr)
                } else{
                    return Promise.reject(error.response.status)
                }
            })
    }

    public static async logout() {
        await Axios.post('auth/logout')
        await Store.dispatch('setLogout')
    }
}

export default AuthAdapter