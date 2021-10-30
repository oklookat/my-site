import type { TGlobalConfig, TRequestConfig, TRequestMethod, THeaders, TRequestBody, TResponse, TError, THook } from './types'

export default class Ducksios {

    public config: TGlobalConfig = {
        timeout: 15000,
        baseURL: null,
        withCredentials: false,
        headers: null,
        hooks: null,
    }

    constructor(config?: TGlobalConfig) {
        if (config) {
            for (const value in config) {
                this.config[value] = config[value]
            }
        }
    }

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
        // send
        xhr.open(method, requestConfig.url, true)
        xhr = this.setRequestHeaders(xhr, requestConfig)
        return this.setupHooksAndSend(xhr, requestConfig)
    }

    private async setupHooksAndSend(xhr: XMLHttpRequest, requestConfig: TRequestConfig): Promise<TResponse> {
        // downloading from server
        xhr.onprogress = (e) => {
            this.onDownload(e, requestConfig)
        }
        // file upload to server
        xhr.upload.onprogress = (e) => {
            this.onUploadProgress(e, requestConfig)
        }
        // file uploaded to server
        xhr.upload.onload = (e) => {
            this.onUploaded(e, requestConfig)
        }
        // return response or error
        return new PromiseWithCancel<TResponse>(xhr, (resolve, reject) => {
            // response from server
            xhr.onload = () => {
                try {
                    const response = this.onResponse(xhr, requestConfig)
                    resolve(response)
                } catch (err) {
                    // HTTP error (in normal cases)
                    reject(err)
                }
            }
            // network error (not HTTP)
            xhr.onerror = () => {
                const err: TError = {
                    type: "network",
                    statusCode: xhr.status,
                    body: null
                }
                this.onError(err, requestConfig)
                reject(err)
            }
            // timeout
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

    // execute hooks when request or response error
    private onError(err: TError, requestConfig: TRequestConfig) {
        const h: THook = {
            name: "onError",
            config: requestConfig,
            data: err
        }
        this.executeHook(h)
    }

    // execute hooks when client send request
    private onRequest(requestConfig: TRequestConfig) {
        const h: THook = {
            name: "onRequest",
            config: requestConfig,
            data: requestConfig
        }
        this.executeHook(h)
    }

    // execute hooks when get response from server
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

    // execute hooks when downloading data from server
    private onDownload(e: ProgressEvent<EventTarget>, requestConfig: TRequestConfig) {
        const h: THook = {
            name: "onDownload",
            config: requestConfig,
            data: e
        }
        this.executeHook(h)
    }

    // execute hooks when upload data to server
    private onUploadProgress(e: ProgressEvent<EventTarget>, requestConfig: TRequestConfig) {
        const h: THook = {
            name: "onUploadProgress",
            config: requestConfig,
            data: e
        }
        this.executeHook(h)
    }

    // execute hooks when data uploaded to server
    private onUploaded(e: ProgressEvent<EventTarget>, requestConfig: TRequestConfig) {
        const h: THook = {
            name: "onUploaded",
            config: requestConfig,
            data: e
        }
        this.executeHook(h)
    }

    // execute hooks depending on global and request config
    private executeHook(h: THook) {
        let hook = this.config.hooks && typeof this.config.hooks[h.name] === 'function'
        if (hook) {
            this.config.hooks[h.name](h.data)
        }
        hook = h.config.hooks && typeof h.config.hooks[h.name] === 'function'
        if (hook) {
            h.config.hooks[h.name](h.data)
        }
    }

    // set url and request params
    private setURL(requestConfig: TRequestConfig): TRequestConfig {
        // set base
        let url = requestConfig.url
        const baseURL = this.config.baseURL
        if (baseURL) {
            url = `${baseURL}/${url}`
        }
        // replace double slashes
        const set = url.match(/([^:]\/{2,3})/g)
        for (const str in set) {
            var replace_with = set[str].substr(0, 1) + '/';
            url = url.replace(set[str], replace_with);
        }
        // set request params
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
        let headers = this.config.headers
        if (headers) {
            set(this.config.headers)
        }
        headers = requestConfig.headers
        if (headers) {
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