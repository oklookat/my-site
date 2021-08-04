import axios  from "axios"
import Store from "@/store/index"
import app from '../../main'

const Axios = axios.create({timeout: 15000})

let token


Axios.defaults.baseURL = process.env.VUE_APP_API_URL
Axios.defaults.headers['Content-Type'] = 'application/json'

Axios.interceptors.request.use(async function (config) {
    // @ts-ignore
    app.$elvenProgress.loadingStart()
    const isAuth = await Store.getters.checkAuth
    if(isAuth){
        token = await Store.getters.getToken
        config.headers.common['Authorization'] = `Elven ${token}`
    }
    return config;
}, function (error) {
    app.$elvenProgress.loadingFinish()
    return Promise.reject(error);
});


Axios.interceptors.response.use(function (response) {
    app.$elvenProgress.loadingFinish()
    return response;
}, function (error) {
    app.$elvenProgress.loadingFinish()
    return Promise.reject(error);
})

export default Axios