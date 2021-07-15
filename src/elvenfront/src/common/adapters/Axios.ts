import axios  from "axios";
import Store from "@/store/index"
const Axios = axios.create({timeout: 15000})


let token


Axios.defaults.baseURL = process.env.VUE_APP_AXIOS_BACKEND_API_URL
Axios.defaults.headers['Content-Type'] = 'application/json'

Axios.interceptors.request.use(async function (config) {
    const isAuth = await Store.getters.checkAuth
    if(isAuth){
        token = await Store.getters.getToken
        config.headers.common['Authorization'] = `Elven ${token}`
    }
    return config;
}, function (error) {

    return Promise.reject(error);
});


Axios.interceptors.response.use(function (response) {

    return response;
}, function (error) {

    return Promise.reject(error);
})

export default Axios