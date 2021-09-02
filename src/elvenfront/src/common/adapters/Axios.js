import axios  from "axios"
import ErrorHandler from "@/common/tools/ErrorHandler";

const Axios = axios.create({timeout: 15000})

const apiURL = import.meta.env.VITE_API_URL


Axios.defaults.baseURL = apiURL
Axios.defaults.headers['Content-Type'] = 'application/json'
Axios.interceptors.request.use(async function (config) {
    window.app.$elvenProgress.loadingStart()
    config.withCredentials = true
    return config;
}, function (error) {
    window.app.$elvenProgress.loadingFinish()
    ErrorHandler.sortError(error)
    return Promise.reject(error);
});


Axios.interceptors.response.use(function (response) {
    app.$elvenProgress.loadingFinish()
    return response;
}, function (error) {
    app.$elvenProgress.loadingFinish()
    ErrorHandler.sortError(error)
    return Promise.reject(error);
})

export default Axios


