import Axios from '@/common/adapters/Axios.js'
import router from "@/router";
import Fetcher from "@/common/adapters/Fetcher";

class AuthAdapter {

    public static async login(username: string, password: string) {
        const data = {
            username: username,
            password: password,
            type: 'cookie',
        }
        return await Axios.post('auth/login', data)
            .then(() => {
                return Promise.resolve()
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
        await router.push({name: 'Login'})
    }

    public static async check(){
        return await Fetcher.check()
            .then((result) =>{
                return Promise.resolve(result)
            })
            .catch((error) =>{
                return Promise.reject(error)
            })
    }
}

export default AuthAdapter