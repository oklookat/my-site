import Duck from '@/duck'
import { AuthStorage } from "@/tools/localStorage";

export default class AuthAdapter {

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
        try {
            await Duck.POST({url: 'auth/logout'})
            AuthStorage.set(false)
        } catch (err) {
            return Promise.reject()
        }
    }
}