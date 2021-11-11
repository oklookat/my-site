import type { TGlobalConfig, TRequestConfig, TRequestMethod, TResponse, TError, THook } from "./types"
import Parser from "./parser"
import Service from "./service"


/** axios-like XHR wrapper */
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

    /** send request (GET) */
    public async GET(config: TRequestConfig): Promise<TResponse> {
        return this.buildAndSend("GET", config)
    }

    /** send request (POST) */
    public async POST(config: TRequestConfig) {
        return this.buildAndSend("POST", config)
    }

    /** send request (PUT) */
    public async PUT(config: TRequestConfig) {
        return this.buildAndSend("PUT", config)
    }

    /** send request (DELETE) */
    public async DELETE(config: TRequestConfig) {
        return this.buildAndSend("DELETE", config)
    }

    /** send request (HEAD) */
    public async HEAD(config: TRequestConfig) {
        return this.buildAndSend("HEAD", config)
    }

    /** send request (OPTIONS) */
    public async OPTIONS(config: TRequestConfig) {
        return this.buildAndSend("OPTIONS", config)
    }

    /** send request (PATCH) */
    public async PATCH(config: TRequestConfig) {
        return this.buildAndSend("PATCH", config)
    }

    /** create XHR, set settings and headers, then send request via {@link setupHooksAndSend} */
    private async buildAndSend(method: TRequestMethod, rc: TRequestConfig): Promise<TResponse> {
        let xhr = new XMLHttpRequest();
        xhr.timeout = this.config.timeout
        // set url
        let url = Service.setBaseURL(rc.url, this.config.baseURL)
        url = Service.setRequestParams(url, rc.params)
        rc.url = url
        // set credentials
        let withCredentials = false
        if (this.config.withCredentials) {
            withCredentials = this.config.withCredentials
        }
        if (rc.withCredentials) {
            withCredentials = rc.withCredentials
        }
        xhr.withCredentials = withCredentials
        // send
        xhr.open(method, rc.url, true)
        xhr = Service.setRequestHeaders(xhr, rc, this.config)
        return this.setupHooksAndSend(xhr, rc)
    }

    /** parse request body, set hooks, send request */
    private async setupHooksAndSend(xhr: XMLHttpRequest, rc: TRequestConfig): Promise<TResponse> {
        const parsed = Parser.requestBody(rc.body, rc, this.config)
        rc = parsed.r
        // when downloading from server
        xhr.onprogress = (e) => {
            this.onDownload(e, rc)
        }
        // when file upload to server
        xhr.upload.onprogress = (e) => {
            this.onUploadProgress(e, rc)
        }
        // when file uploaded to server
        xhr.upload.onload = (e) => {
            this.onUploaded(e, rc)
        }
        // return response or error
        return new Promise((resolve, reject) => {
            // when response from server
            xhr.onload = () => {
                try {
                    const response = this.onResponse(xhr, rc)
                    resolve(response)
                } catch (err) {
                    // HTTP error (in normal cases)
                    reject(err)
                }
            }
            // when network error (not HTTP)
            xhr.onerror = () => {
                const err: TError = {
                    type: "network"
                }
                this.onError(err, rc)
                reject(err)
            }
            // when timeout
            xhr.ontimeout = () => {
                const err: TError = {
                    type: "timeout"
                }
                this.onError(err, rc)
                reject(err)
            }
            // is cancel token provided
            if (rc.cancelToken) {
                // rewrite CancelToken.cancel function
                rc.cancelToken.cancel = (message?: string) => {
                    // now when calls cancel(), executes this:
                    xhr.abort()
                    const err: TError = {
                        type: "cancel",
                        message: message
                    }
                    this.onError(err, rc)
                    reject(err)
                }
            }
            xhr.send(parsed.b)
            this.onRequest(rc)
        })
    }

    /** execute hooks when request or response error */
    private onError(err: TError, rc: TRequestConfig) {
        const h: THook = {
            name: "onError",
            config: rc,
            data: err
        }
        this.executeHooks(h)
    }

    /** execute hooks when client send request */
    private onRequest(rc: TRequestConfig) {
        const h: THook = {
            name: "onRequest",
            config: rc,
            data: rc
        }
        this.executeHooks(h)
    }

    /** execute hooks when get response from server */
    private onResponse(xhr: XMLHttpRequest, rc: TRequestConfig): TResponse {
        const statusCode = xhr.status
        const body = Parser.responseBody(xhr.response)
        if (statusCode != 200) {
            const err: TError = {
                type: "response",
                statusCode: statusCode,
                body: body
            }
            // execute user hooks
            this.onError(err, rc)
            // send error
            throw err
        }
        const resp: TResponse = {
            body: body,
            statusCode: statusCode
        }
        // execute user hooks
        const h: THook = {
            name: "onResponse",
            config: rc,
            data: resp
        }
        this.executeHooks(h)
        // send response
        return resp
    }

    /** execute hooks when downloading data from server */
    private onDownload(e: ProgressEvent<EventTarget>, rc: TRequestConfig) {
        const h: THook = {
            name: "onDownload",
            config: rc,
            data: e
        }
        this.executeHooks(h)
    }

    /** execute hooks when upload data to server */
    private onUploadProgress(e: ProgressEvent<EventTarget>, rc: TRequestConfig) {
        const h: THook = {
            name: "onUploadProgress",
            config: rc,
            data: e
        }
        this.executeHooks(h)
    }

    /** execute hooks when data uploaded to server */
    private onUploaded(e: ProgressEvent<EventTarget>, rc: TRequestConfig) {
        const h: THook = {
            name: "onUploaded",
            config: rc,
            data: e
        }
        this.executeHooks(h)
    }

    /** execute hooks depending on global and request config */
    private executeHooks(h: THook) {
        let hook = this.config.hooks && h.name in this.config.hooks && typeof this.config.hooks[h.name] === 'function'
        if (hook) {
            this.config.hooks[h.name](h.data as any)
        }
        hook = h.config.hooks && h.name in h.config.hooks && typeof h.config.hooks[h.name] === 'function'
        if (hook) {
            h.config.hooks[h.name](h.data as any)
        }
    }
}