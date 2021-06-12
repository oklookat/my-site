import Axios from '@/common/adapters/Axios'
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
                return Promise.reject('токен не получен.')
            })
            .catch(error =>{
                const readableErr = error.response.data.error
                if(readableErr){
                    return Promise.reject(readableErr)
                }
                return Promise.reject(error.response.status)
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