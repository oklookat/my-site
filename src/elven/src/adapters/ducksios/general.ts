import type { TGlobalConfig, TRequestConfig, TRequestMethod, THeaders, TRequestBody, TResponse, TError, THook } from './types'

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
        // call hooks when downloading from server
        xhr.onprogress = (e) => {
            this.onDownload(e, requestConfig)
        }
        // call hooks when upload to server
        xhr.upload.onprogress = (e) => {
            this.onUploadProgress(e, requestConfig)
        }
        // call hooks when data uploaded to server
        xhr.upload.onload = (e) => {
            this.onUploaded(e, requestConfig)
        }
        // return response or error
        return new Promise((resolve, reject) => {
            // call hooks when we get response from server
            xhr.onload = () => {
                try {
                    const response = this.onResponse(xhr, requestConfig)
                    resolve(response)
                } catch (err) {
                    // HTTP error (mostly)
                    reject(err)
                }
            }
            // call hooks when network error (not HTTP)
            xhr.onerror = () => {
                const err: TError = {
                    type: "network",
                    statusCode: xhr.status,
                    body: null
                }
                this.onError(err, requestConfig)
                reject(err)
            }
            // call hooks when timeout
            xhr.ontimeout = () => {
                const err: TError = {
                    type: "timeout",
                    statusCode: xhr.status,
                    body: null
                }
                this.onError(err, requestConfig)
                reject(err)
            }
            xhr.send(this.parseRequestBody(requestConfig.body))
            this.onRequest(requestConfig)
        })
    }

    //////////// hooks

    private executeUserHooks(name: THook, requestConfig: TRequestConfig, data: TRequestConfig | TError | ProgressEvent<EventTarget>) {
        let hook = this.config.hooks && typeof this.config.hooks[name] === 'function'
        if (hook) {
            this.config.hooks[name](data)
        }
        hook = requestConfig.hooks && typeof requestConfig.hooks[name] === 'function'
        if (hook) {
            requestConfig.hooks[name](data)
        }
    }

    // when error
    private onError(err: TError, requestConfig: TRequestConfig) {
        this.executeUserHooks("onError", requestConfig, err)

    }

    // when get response from server
    private onResponse(xhr: XMLHttpRequest, requestConfig: TRequestConfig): TResponse {
        const statusCode = xhr.status
        const body = this.parseResponseBody(xhr)
        if (statusCode != 200) {
            const err: TError = {
                type: "response",
                statusCode: statusCode,
                body: body
            }
            // execute user hooks
            this.onError(err, requestConfig)
            // send error
            throw err
        }
        const resp: TResponse = {
            body: body,
            statusCode: statusCode
        }
        // execute user hooks
        let hook = this.config.hooks && typeof this.config.hooks.onResponse === 'function'
        if (hook) {
            this.config.hooks.onResponse(resp)
        }
        hook = requestConfig.hooks && typeof requestConfig.hooks.onResponse === 'function'
        if (hook) {
            requestConfig.hooks.onResponse(resp)
        }
        // send response
        return resp
    }

    // execute user hooks when client send request
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

    // execute user hooks when downloading data from server
    private onDownload(e: ProgressEvent<EventTarget>, requestConfig: TRequestConfig) {
        let hook = this.config.hooks && typeof this.config.hooks.onDownload === 'function'
        if (hook) {
            this.config.hooks.onDownload(e)
        }
        hook = requestConfig.hooks && typeof requestConfig.hooks.onDownload === 'function'
        if (hook) {
            requestConfig.hooks.onDownload(e)
        }
    }

    // execute user hooks when upload data to server
    private onUploadProgress(e: ProgressEvent<EventTarget>, requestConfig: TRequestConfig) {
        let hook = this.config.hooks && typeof this.config.hooks.onUploadProgress === 'function'
        if (hook) {
            this.config.hooks.onUploadProgress(e)
        }
        hook = requestConfig.hooks && typeof requestConfig.hooks.onUploadProgress === 'function'
        if (hook) {
            requestConfig.hooks.onUploadProgress(e)
        }
    }

    // execute user hooks when data uploaded to server
    private onUploaded(e: ProgressEvent<EventTarget>, requestConfig: TRequestConfig) {
        let hook = this.config.hooks && typeof this.config.hooks.onUploaded === 'function'
        if (hook) {
            this.config.hooks.onUploaded(e)
        }
        hook = requestConfig.hooks && typeof requestConfig.hooks.onUploaded === 'function'
        if (hook) {
            requestConfig.hooks.onUploaded(e)
        }
    }

    //////////// parsers

    // set url and request params
    private setURL(requestConfig: TRequestConfig): TRequestConfig {
        // set base
        let url = requestConfig.url
        const baseURL = this.config.baseURL
        if (baseURL) {
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
        if (!body) {
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
}