import type { GlobalConfig, RequestConfig, Response, RequestFail, HookOutput } from "./types"
import { RequestMethod, RequestError, HookName } from "./types"
import Parse from "./utils/parse"
import Service from "./utils/service"
import { Validator } from "./utils/validator"


/** axios-like XHR wrapper */
export default class Ducksios {

    public config: GlobalConfig = {
        timeout: 15000,
        baseURL: null,
        withCredentials: false,
        headers: null,
        hooks: null,
    }

    constructor(config?: GlobalConfig) {
        if (!config) {
            return
        }
        for (const value in config) {
            this.config[value] = config[value]
        }
    }

    /** send request (GET) */
    public async GET(config: RequestConfig): Promise<Response> {
        return this.buildAndSend(RequestMethod.GET, config)
    }

    /** send request (POST) */
    public async POST(config: RequestConfig) {
        return this.buildAndSend(RequestMethod.POST, config)
    }

    /** send request (PUT) */
    public async PUT(config: RequestConfig) {
        return this.buildAndSend(RequestMethod.PUT, config)
    }

    /** send request (DELETE) */
    public async DELETE(config: RequestConfig) {
        return this.buildAndSend(RequestMethod.DELETE, config)
    }

    /** send request (HEAD) */
    public async HEAD(config: RequestConfig) {
        return this.buildAndSend(RequestMethod.HEAD, config)
    }

    /** send request (OPTIONS) */
    public async OPTIONS(config: RequestConfig) {
        return this.buildAndSend(RequestMethod.OPTIONS, config)
    }

    /** send request (PATCH) */
    public async PATCH(config: RequestConfig) {
        return this.buildAndSend(RequestMethod.PATCH, config)
    }

    /** create XHR, set settings and headers, then send request via {@link setupHooksAndSend} */
    private async buildAndSend(method: RequestMethod, rc: RequestConfig): Promise<Response> {
        let xhr = new XMLHttpRequest();
        xhr.timeout = this.config.timeout
        // url
        let url = Service.setBaseURL(rc.url, this.config.baseURL)
        url = Service.setRequestParams(url, rc.params)
        rc.url = url
        // credentials
        let withCredentials = false
        if (rc.withCredentials) {
            withCredentials = rc.withCredentials
        } else if (this.config.withCredentials) {
            withCredentials = this.config.withCredentials
        }
        xhr.withCredentials = withCredentials
        // send
        xhr.open(method, rc.url, true)
        Service.setRequestHeaders(xhr, rc, this.config)
        return this.setupHooksAndSend(xhr, rc)
    }

    /** set config, parse request body, set hooks, send request */
    private async setupHooksAndSend(xhr: XMLHttpRequest, rc: RequestConfig): Promise<Response> {
        const body = Parse.requestBody(rc.body, rc, this.config)
        // return response or error
        return new Promise((resolve, reject) => {
            // when downloading from server
            xhr.onprogress = (e) => {
                this.onDownload(e, rc)
            }
            // when file uploading to server
            xhr.upload.onprogress = (e) => {
                this.onUploadProgress(e, rc)
            }
            // when file uploaded to server
            xhr.upload.onload = (e) => {
                this.onUploaded(e, rc)
            }
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
                const err: RequestFail = {
                    type: RequestError.network
                }
                this.onError(err, rc)
                reject(err)
            }
            // when timeout
            xhr.ontimeout = () => {
                const err: RequestFail = {
                    type: RequestError.timeout
                }
                this.onError(err, rc)
                reject(err)
            }
            // if cancel token provided
            if (rc.cancelToken) {
                // rewrite CancelToken.cancel function
                rc.cancelToken.cancel = (message?: string) => {
                    // executes when CancelToken.cancel() was called
                    xhr.abort()
                    const err: RequestFail = {
                        type: RequestError.cancel,
                        message: message
                    }
                    this.onError(err, rc)
                    reject(err)
                }
            }
            xhr.send(body)
            this.onRequest(rc)
        })
    }

    /** execute hooks when request or response error */
    private onError(err: RequestFail, rc: RequestConfig) {
        const h: HookOutput = {
            name: HookName.onError,
            config: rc,
            data: err
        }
        this.executeHooks(h)
    }

    /** execute hooks when client send request */
    private onRequest(rc: RequestConfig) {
        const h: HookOutput = {
            name: HookName.onRequest,
            config: rc,
            data: rc
        }
        this.executeHooks(h)
    }

    /** execute hooks when get response from server */
    private onResponse(xhr: XMLHttpRequest, rc: RequestConfig): Response {
        const statusCode = xhr.status
        const body = Parse.responseBody(xhr.response)
        const is2xx = (statusCode - 200) >= 0 && (statusCode - 200) <= 99
        // if statusCode not 2**
        if (!is2xx) {
            const err: RequestFail = {
                type: RequestError.response,
                statusCode: statusCode,
                body: body
            }
            // execute user hooks
            this.onError(err, rc)
            // send error
            throw err
        }
        const resp: Response = {
            body: body,
            statusCode: statusCode
        }
        // execute user hooks
        const h: HookOutput = {
            name: HookName.onResponse,
            config: rc,
            data: resp
        }
        this.executeHooks(h)
        // send response
        return resp
    }

    /** execute hooks when downloading data from server */
    private onDownload(e: ProgressEvent<EventTarget>, rc: RequestConfig) {
        const h: HookOutput = {
            name: HookName.onDownload,
            config: rc,
            data: e
        }
        this.executeHooks(h)
    }

    /** execute hooks when upload data to server */
    private onUploadProgress(e: ProgressEvent<EventTarget>, rc: RequestConfig) {
        const h: HookOutput = {
            name: HookName.onUploadProgress,
            config: rc,
            data: e
        }
        this.executeHooks(h)
    }

    /** execute hooks when data uploaded to server */
    private onUploaded(e: ProgressEvent<EventTarget>, rc: RequestConfig) {
        const h: HookOutput = {
            name: HookName.onUploaded,
            config: rc,
            data: e
        }
        this.executeHooks(h)
    }

    /** execute hooks depending on global and request config */
    private executeHooks(h: HookOutput) {
        const globalHookAvailable = Validator.isHookAvailable(h.name, this.config.hooks)
        if (globalHookAvailable) {
            this.config.hooks[h.name](h as any)
        }
        const localHookAvailable = Validator.isHookAvailable(h.name, h.config.hooks)
        if (localHookAvailable) {
            h.config.hooks[h.name](h as any)
        }
    }
}