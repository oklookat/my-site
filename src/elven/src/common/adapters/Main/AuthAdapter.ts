import { Axios } from '@/common/adapters/Axios'
import { AuthStorage } from "@/common/tools/LStorage";

class AuthAdapter {

    public static async login(username: string, password: string) {
        const data = {
            username: username,
            password: password,
            type: 'cookie',
        }
        try {
            await Axios.post('auth/login', data)
            AuthStorage.set(true)
            return Promise.resolve()
        } catch {
            return Promise.reject()
        }
    }

    public static async logout() {
        AuthStorage.set(false)
        await Axios.post('auth/logout').catch(() => {return Promise.reject()})
    }
}

export default AuthAdapter