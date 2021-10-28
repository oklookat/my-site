import type { TGlobalConfig, TRequestConfig, TRequestMethod, THeaders, TRequestBody, TResponse, TError } from './types'

export default class Ducksios {

    public config: TGlobalConfig = {
        timeout: 15000,
        baseURL: null,
        withCredentials: false,
        headers: null,
        hooks: null,
        baseAuth: false,
    }

    constructor(config?: TGlobalConfig) {
        if (config) {
            for (const value in config) {
                this.config[value] = config[value]
            }
        }
    }

    //////////// methods

    public async GET(config: TRequestConfig) {
        return this.buildAndSend("GET", config)
    }

    public async POST(config: TRequestConfig) {
        return this.buildAndSend("POST", config)
    }

    public async PUT(config: TRequestConfig) {
        return this.buildAndSend("PUT", config)
    }

    public async DELETE(config: TRequestConfig) {
        return this.buildAndSend("DELETE", config)
    }

    public async HEAD(config: TRequestConfig) {
        return this.buildAndSend("HEAD", config)
    }

    public async OPTIONS(config: TRequestConfig) {
        return this.buildAndSend("OPTIONS", config)
    }

    public async PATCH(config: TRequestConfig) {
        return this.buildAndSend("PATCH", config)
    }

    //////////// senders

    private async buildAndSend(method: TRequestMethod, requestConfig: TRequestConfig): Promise<TResponse> {
        let xhr = new XMLHttpRequest();
        xhr.timeout = this.config.timeout
        requestConfig = this.setURL(requestConfig)
        // credentials
        let withCredentials = false
        if (this.config.withCredentials) {
            withCredentials = this.config.withCredentials
        }
        if (requestConfig.withCredentials) {
            withCredentials = requestConfig.withCredentials
        }
        xhr.withCredentials = withCredentials
        // base auth
        let baseAuthUser: string | undefined = undefined
        let baseAuthPassword: string | undefined = undefined
        if (this.config.baseAuth || requestConfig.baseAuth) {
            let user = this.config.baseAuthUser
            if (user) {
                baseAuthUser = user
            }
            user = requestConfig.baseAuthUser
            if (user) {
                baseAuthUser = user
            }
            let password = this.config.baseAuthPassword
            if (password) {
                baseAuthPassword = password
            }
            password = requestConfig.baseAuthPassword
            if (password) {
                baseAuthPassword = password
            }
        }
        // send
        const url = requestConfig.url
        xhr.open(method, url, true, baseAuthUser, baseAuthPassword)
        xhr = this.setRequestHeaders(xhr, requestConfig)
        return this.send(xhr, requestConfig)
    }

    // send request and setup hooks
    private async send(xhr: XMLHttpRequest, requestConfig: TRequestConfig): Promise<TResponse> {
        let body = this.parseRequestBody(requestConfig.body)
        xhr.send(body)
        // after send request
        this.onRequest(requestConfig)
        return new Promise((resolve, reject) => {
            // when we get response from server
            xhr.onload = () => {
                this.onResponse(xhr, requestConfig)
                    .then((response) => {
                        resolve(response)
                    })
                    .catch(err => {
                        reject(err as TError)
                    })
            }
            // network error (not HTTP)
            xhr.onerror = () => {
                const err: TError = {
                    type: "network",
                    statusCode: 0,
                    body: null
                }
                reject(this.onRequestError(requestConfig, err))
            }
            // timeout
            xhr.ontimeout = () => {
                const err: TError = {
                    type: "timeout",
                    statusCode: 408,
                    body: null
                }
                reject(this.onResponseError(requestConfig, err))
            }
        })
    }

    //////////// parsers

    // set url and request params
    private setURL(requestConfig: TRequestConfig): TRequestConfig {
        // set base
        let url = requestConfig.url
        const baseURL = this.config.baseURL
        if(baseURL) {
            url = `${baseURL}/${url}`
        }
        // control slashes
        const set = url.match(/([^:]\/{2,3})/g)
        for (const str in set) {
            var replace_with = set[str].substr(0, 1) + '/';
            url = url.replace(set[str], replace_with);
        }
        // set params
        let urlObj = new URL(url)
        if (requestConfig.params) {
            for (const param in requestConfig.params) {
                urlObj.searchParams.set(param, requestConfig.params[param].toString())
            }
        }
        requestConfig.url = urlObj.toString()
        return requestConfig
    }

    // set request headers by request and global config
    private setRequestHeaders(xhr: XMLHttpRequest, requestConfig: TRequestConfig): XMLHttpRequest {
        const set = (headers: THeaders) => {
            for (const header in headers) {
                xhr.setRequestHeader(header, headers[header].toString())
            }
        }
        const globalHeaders = this.config.headers
        if (globalHeaders) {
            set(this.config.headers)
        }
        const localHeaders = requestConfig.headers
        if (localHeaders) {
            set(requestConfig.headers)
        }
        return xhr
    }

    // TODO: auto parse json and add request cancel, progress

    // parse body to json (if object) before send request. Returns body.
    private parseRequestBody(body: TRequestBody): any {
        if(!body) {
            return
        }
        // const toJSON = !(body instanceof Blob) && !(body instanceof Buffer) && !(body instanceof FormData) 
        // && !(body instanceof URLSearchParams) && !(body instanceof ReadableStream) && typeof body === 'object'
        // if (toJSON) {
        //     try {
        //         body = JSON.stringify(body)
        //     } catch (err) {

        //     }
        // } else if (typeof body === 'number') {
        //     body = body.toString()
        // }
        return body
    }

    // try to parse body. Returns body.
    private parseResponseBody(xhr: XMLHttpRequest): any {
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

    //////////// hooks

    // when get response from server
    private async onResponse(xhr: XMLHttpRequest, requestConfig: TRequestConfig): Promise<TResponse> {
        return new Promise((resolve, reject) => {
            const statusCode = xhr.status
            const body = this.parseResponseBody(xhr)
            if (statusCode != 200) {
                const err: TError = {
                    type: "response",
                    statusCode: statusCode,
                    body: body
                }
                reject(this.onResponseError(requestConfig, err))
                return
            }
            const resp: TResponse = {
                body: body,
                statusCode: statusCode
            }
            let hook = this.config.hooks && typeof this.config.hooks.onResponse === 'function'
            if (hook) {
                this.config.hooks.onResponse(resp)
            }
            hook = requestConfig.hooks && typeof requestConfig.hooks.onResponse === 'function'
            if (hook) {
                requestConfig.hooks.onResponse(resp)
            }
            resolve(resp)
        })
    }

    // executes user hooks when client send request
    private onRequest(requestConfig: TRequestConfig) {
        let hook = this.config.hooks && typeof this.config.hooks.onRequest === 'function'
        if (hook) {
            this.config.hooks.onRequest(requestConfig)
        }
        hook = requestConfig.hooks && typeof requestConfig.hooks.onRequest === 'function'
        if (hook) {
            requestConfig.hooks.onRequest(requestConfig)
        }
    }

    // executes user hooks when error like no internet etc
    private onRequestError(requestConfig: TRequestConfig, err: TError): TError {
        let hook = this.config.hooks && typeof this.config.hooks.onRequestError === 'function'
        if (hook) {
            this.config.hooks.onRequestError(err)
        }
        hook = requestConfig.hooks && typeof requestConfig.hooks.onRequestError === 'function'
        if (hook) {
            requestConfig.hooks.onRequestError(err)
        }
        return err
    }

    // executes user hooks when HTTP error like 404
    private onResponseError(requestConfig: TRequestConfig, err: TError): TError {
        let hook = this.config.hooks && typeof this.config.hooks.onResponseError === 'function'
        if (hook) {
            this.config.hooks.onResponseError(err)
        }
        hook = requestConfig.hooks && typeof requestConfig.hooks.onResponseError === 'function'
        if (hook) {
            requestConfig.hooks.onResponseError(err)
        }
        return err
    }
}