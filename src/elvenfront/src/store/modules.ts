import router from "@/router";

interface IAuthInterface {
    isLoggedIn: boolean
    token: string | null
}


export const authStore = {
    state: {
        auth: {
            isLoggedIn: false,
            token: null
        }
    },
    mutations: {
        setAuth(state, auth: IAuthInterface) {
            state.auth = auth
        },
    },
    actions: {
        async setToken({commit}, token: string) {
            const auth = {
                isLoggedIn: true,
                token: token
            }
            commit('setAuth', auth)
            await router.push({name: 'Index'})
        },
        async setLogout({commit}){
            const auth = {
                isLoggedIn: false,
                token: null
            }
            commit('setAuth', auth)
            await router.push({name: 'Login'})
        },
    },
    getters: {
        async getToken(state) {
            return Promise.resolve(state.auth.token)
        },
        async checkAuth(state){
            const auth = state.auth
            if(auth.isLoggedIn && auth.token){
                return Promise.resolve(true)
            }
            return Promise.resolve(false)
        }
    }
}