import { Axioser } from "./Axioser"
import type { TAxiosParams, TRequestInterceptor, TResponseInterceptor } from "./Axioser"
import ErrorHandler from "@/common/tools/ErrorHandler"
const apiURL = import.meta.env.VITE_API_URL

/* -------- MAIN -------- */
const mainRequest: TRequestInterceptor = {
    onFulfilled: (config) => {
        window.$elvenProgress.start()
        config.withCredentials = true
        return config;
    },
    onRejected: (err) => {
        window.$elvenProgress.finish()
        ErrorHandler(err)
        return Promise.reject(err);
    }
}

const mainResponse: TResponseInterceptor = {
    onFulfilled: (response) => {
        window.$elvenProgress.finish()
        return response;
    },
    onRejected: (err) => {
        window.$elvenProgress.finish()
        ErrorHandler(err)
        return Promise.reject(err);
    }
}

const mainParams: TAxiosParams = {
    apiURL: apiURL as string,
    requestInterceptor: mainRequest,
    responseInterceptor: mainResponse
}

const mainAxios = new Axioser(mainParams)
export const Axios = mainAxios.instance