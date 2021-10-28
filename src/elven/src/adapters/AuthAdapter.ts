import Duck from './Duck'
import { AuthStorage } from "@/tools/LocalStorage";

class AuthAdapter {

    public static async login(username: string, password: string) {
        const data = {
            username: username,
            password: password,
            type: 'cookie',
        }
        try {
            await Duck.POST({url: 'auth/login', body: data})
            AuthStorage.set(true)
            return Promise.resolve()
        } catch {
            return Promise.reject()
        }
    }

    public static async logout() {
        AuthStorage.set(false)
        try {
            await Duck.POST({url: 'auth/logout'})
        } catch (err) {
            return Promise.reject()
        }
    }
}

export default AuthAdapter