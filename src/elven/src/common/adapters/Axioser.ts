import axios, { AxiosRequestConfig, AxiosResponse } from "axios"
// TS wrapper for Axios

export type TRequestInterceptor = {
    onFulfilled?: (config: AxiosRequestConfig<any>) => AxiosRequestConfig<any> | Promise<AxiosRequestConfig<any>>
    onRejected?: (error: any) => any
}


export type TResponseInterceptor = {
    onFulfilled?: (response: AxiosResponse<unknown, any>) => AxiosResponse<unknown, any> | Promise<AxiosResponse<unknown, any>>
    onRejected?: (error: any) => any
}


export type TAxiosParams = {
    apiURL: string
    requestInterceptor?: TRequestInterceptor
    responseInterceptor?: TResponseInterceptor
}


export class Axioser {

    public instance = axios.create({ timeout: 15000 })

    constructor(params: TAxiosParams) {
        this.instance.defaults.baseURL = params.apiURL
        this.instance.defaults.headers['Content-Type'] = 'application/json'
        this.initInterceptors(params.requestInterceptor, params.responseInterceptor)
    }

    private initInterceptors(request?: TRequestInterceptor, response?: TResponseInterceptor) {
        let requestFull = undefined
        let requestRej = undefined
        let responseFull = undefined
        let responseRej = undefined
        if (request) {
            if (request.onFulfilled) {
                requestFull = request.onFulfilled
            }
            if (request.onRejected) {
                requestRej = request.onRejected
            }
        }
        if (response) {
            if (response.onFulfilled) {
                responseFull = response.onFulfilled
            }
            if (response.onRejected) {
                responseRej = response.onRejected
            }
        }
        this.instance.interceptors.request.use(requestFull, requestRej)
        this.instance.interceptors.response.use(responseFull, responseRej)
    }
}


