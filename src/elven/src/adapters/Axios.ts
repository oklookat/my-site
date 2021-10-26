import { Axioser } from "./Axioser"
import type { TAxiosParams, TRequestInterceptor, TResponseInterceptor } from "./Axioser"
import { AdapterError } from "@/tools/ErrorHandler"
import { Env } from "@/tools/Paths"

const apiURL = Env.getAPI()

// TODO: try use fetch and make interceptors in them(?)
const mainRequest: TRequestInterceptor = {
    onFulfilled: (config) => {
        window.$elvenProgress.start()
        config.withCredentials = true
        return config;
    },
    onRejected: (err) => {
        window.$elvenProgress.finish()
        AdapterError.handle(err)
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
        AdapterError.handle(err)
        return Promise.reject(err);
    }
}

const mainParams: TAxiosParams = {
    apiURL: apiURL,
    requestInterceptor: mainRequest,
    responseInterceptor: mainResponse
}

const mainAxios = new Axioser(mainParams)
export const Axios = mainAxios.instance