import { IGlobalConfig, TRequestMethod, IRequestConfig, THeader, IError, TResponse } from './types'

class Ducksios {

    public config: IGlobalConfig = {
        timeout: 15000,
        baseAuth: false
    }

    constructor(config?: IGlobalConfig) {
        if (config) {
            if (config.timeout) {
                this.config.timeout = config.timeout
            }
            if (config.baseURL) {
                this.config.baseURL = config.baseURL
            }
            if (config.headers) {
                this.config.headers = config.headers
            }
        }
    }

    public async GET(config: IRequestConfig) {
        this.buildXHR("GET", config)
    }

    public async POST(config: IRequestConfig) {
        this.buildXHR("POST", config)
    }

    public async PUT(config: IRequestConfig) {
        this.buildXHR("PUT", config)
    }

    public async DELETE(config: IRequestConfig) {
        this.buildXHR("DELETE", config)
    }

    public async HEAD(config: IRequestConfig) {
        this.buildXHR("HEAD", config)
    }

    public async OPTIONS(config: IRequestConfig) {
        this.buildXHR("OPTIONS", config)
    }

    public async PATCH(config: IRequestConfig) {
        this.buildXHR("PATCH", config)
    }

    // cancel request
    public stop() {

    }

    // TODO: array with xhr(?) and request config, and add on load end (destroy). And add onprogress.

    // create XMLHttpRequest (this.xhr)
    private buildXHR(method: TRequestMethod, requestConfig: IRequestConfig): XMLHttpRequest {
        let xhr = new XMLHttpRequest();
        xhr.timeout = this.config.timeout
        // url + baseURL
        let url = requestConfig.url
        if (this.config.baseURL) {
            url = `${this.config.baseURL}/${url}`
        }
        // credentials
        let withCredentials = false
        if (this.config.withCredentials) {
            withCredentials = this.config.withCredentials
        } else if (requestConfig.withCredentials) {
            withCredentials = requestConfig.withCredentials
        }
        xhr.withCredentials = withCredentials
        // base auth
        let baseAuthUser: string | undefined = undefined
        let baseAuthPassword: string | undefined = undefined
        if (this.config.baseAuth || requestConfig.baseAuth) {
            if (this.config.baseAuthUser) {
                baseAuthUser = this.config.baseAuthUser
            }
            if (requestConfig.baseAuthUser) {
                baseAuthUser = requestConfig.baseAuthUser
            }
            if (this.config.baseAuthPassword) {
                baseAuthPassword = this.config.baseAuthPassword
            }
            if (requestConfig.baseAuthPassword) {
                baseAuthPassword = requestConfig.baseAuthPassword
            }
        }
        // hooks
        xhr = this.bootHooks(xhr, requestConfig)
        // open
        xhr.open(method, url, true, baseAuthUser, baseAuthPassword)
        xhr = this.parseHeaders(xhr, requestConfig)
        return xhr
    }

    // set headers by request and global config
    private parseHeaders(xhr: XMLHttpRequest, requestConfig: IRequestConfig): XMLHttpRequest {
        const set = (headers: Array<THeader>) => {
            for (const header of headers) {
                xhr.setRequestHeader(header.name, header.value)
            }
        }
        const globalHeaders = this.config.headers && this.config.headers.length > 0
        if (globalHeaders) {
            set(this.config.headers)
        }
        const localHeaders = requestConfig.headers && requestConfig.headers.length > 0
        if (localHeaders) {
            set(requestConfig.headers)
        }
        return xhr
    }

    // try to parse body if object. Returns body.
    private parseBody(xhr: XMLHttpRequest): any {
        let body = xhr.response
        if (body) {
            try {
                body = JSON.parse(body)
            } catch (err) {
            }
        } else {
            body = null
        }
        return body
    }

    // setup hooks by request and global config
    private bootHooks(xhr: XMLHttpRequest, requestConfig: IRequestConfig): XMLHttpRequest {
        // network error (not HTTP)
        xhr.onerror = () => {
            const err: IError = {
                type: "network",
                statusCode: 0,
                body: null
            }
            this.onRequestError(requestConfig, err)
        }
        // timeout
        xhr.ontimeout = () => {
            const err: IError = {
                type: "timeout",
                statusCode: 408,
                body: null
            }
            this.onResponseError(requestConfig, err)
        }
        // on request send
        xhr.onloadstart = () => {
            this.onRequest(requestConfig)
        }
        // on response
        xhr.onload = () => {
            this.onResponse(xhr, requestConfig)
        }
        return xhr
    }

    // hook: when client send request
    private onRequest(requestConfig: IRequestConfig) {
        const globalRequestHook = this.config.hooks && typeof this.config.hooks.onRequest === 'function'
        if (globalRequestHook) {
            this.config.hooks.onRequest(requestConfig)
        }
        const localRequestHook = requestConfig.hooks && typeof requestConfig.hooks.onRequest === 'function'
        if (localRequestHook) {
            requestConfig.hooks.onRequest(requestConfig)
        }
    }

    // hook: when server send reponse
    private onResponse(xhr: XMLHttpRequest, requestConfig: IRequestConfig) {
        const statusCode = xhr.status
        const body = this.parseBody(xhr)
        if (statusCode != 200) {
            const err: IError = {
                type: "response",
                statusCode: statusCode,
                body: body
            }
            this.onResponseError(requestConfig, err)
            return
        }
        const resp: TResponse = {
            body: body,
            statusCode: statusCode
        }
    }

    // hook: when error like no internet etc
    private onRequestError(requestConfig: IRequestConfig, err: IError) {
        const globalRequestErrorHook = this.config.hooks && typeof this.config.hooks.onRequestError === 'function'
        if (globalRequestErrorHook) {
            this.config.hooks.onRequestError(err)
        }
        const localRequestErrorHook = requestConfig.hooks && typeof requestConfig.hooks.onRequestError === 'function'
        if (localRequestErrorHook) {
            requestConfig.hooks.onRequestError(err)
        }
    }

    // hook: when HTTP error like 404
    private onResponseError(requestConfig: IRequestConfig, err: IError) {
        const globalResponseErrorHook = this.config.hooks && typeof this.config.hooks.onResponseError === 'function'
        if (globalResponseErrorHook) {
            this.config.hooks.onResponseError(err)
        }
        const localResponseErrorHook = requestConfig.hooks && typeof requestConfig.hooks.onResponseError === 'function'
        if (localResponseErrorHook) {
            requestConfig.hooks.onResponseError(err)
        }
    }

    // hook: when XHR end
    private onLoadEnd() {

    }
}