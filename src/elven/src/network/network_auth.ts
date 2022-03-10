import Duck from '@/network'
import type { body } from '@/types/auth';
import { AuthStorage } from "@/tools/storage"

export default class NetworkAuth {

    public static async login(username: string, password: string) {
        const data: body = {
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