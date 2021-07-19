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
                if(token){
                    await Store.dispatch('setToken', token)
                    return Promise.resolve()
                }
                return Promise.reject('Токен не получен.')
            })
            .catch(error =>{
                if(!error.response){
                    return Promise.reject('Произошла ошибка. Попробуйте позже.')
                }
                const readableErr = error.response.status
                return Promise.reject(readableErr)
            })
    }

    public static async logout() {
        const token = await Store.getters.getToken
        Axios.defaults.headers['Authorization'] = `Elven ${token}`
        await Store.dispatch('setLogout')
        await Axios.post('auth/logout')
    }
}

export default AuthAdapter