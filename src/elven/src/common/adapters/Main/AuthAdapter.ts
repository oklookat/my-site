import Axios from '@/common/adapters/Axios.js'
import {AuthStorage} from "@/common/tools/LStorage";

class AuthAdapter {

    public static async login(username: string, password: string) {
        const data = {
            username: username,
            password: password,
            type: 'cookie',
        }
        try {
            await Axios.post('auth/login', data)
            // set auth state after request, because we can get error when username and password not valid
            AuthStorage.set(true)
            return Promise.resolve()
        } catch (err){
        }
    }

    public static async logout() {
        AuthStorage.set(false)
        await Axios.post('auth/logout').catch(() => {})
        await router.push({name: 'Login'})
    }
}

export default AuthAdapter