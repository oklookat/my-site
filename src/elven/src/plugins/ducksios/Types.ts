import type CancelToken from "./cancel"

/** request method */
export type TRequestMethod = "GET" | "POST" | "PUT" | "DELETE" | "HEAD" | "OPTIONS" | "PATCH"

/** available request body */
export type TRequestBody = string | number | object | Blob | BufferSource | FormData | URLSearchParams | ReadableStream

/** typical request headers */
export type THeaders = {
    [name: string]: string | number
}

/** typical request params */
export type TRequestParams = {
    [name: string]: string | number | boolean
}

/** response from server */
export type TResponse = {
    body: any
    statusCode: number
}

/** represents request errors */
export type TError = {
    type: "timeout" | "network" | "request"
} | TResponse & {
    type: "response"
} | {
    type: "cancel"
    message?: string | number
}

/** global & local configs extends this */
export type TBasicConfig = {
    /** @see https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest/withCredentials */
    withCredentials?: boolean
    headers?: THeaders
    hooks?: IHooks
}

/** global configuration for all requests */
export type TGlobalConfig = TBasicConfig & {
    /** how long to wait for a response from the server (in ms)  */
    timeout?: number
    /**
     * URL like: https://example.com.
     * Any request will paste this url before request url.
     */
    baseURL?: string
}

/** configuration for one request */
export type TRequestConfig = TBasicConfig & {
    /** URL or path like '/hello/world' if {@link TGlobalConfig.baseURL baseURL} setted */
    url: string
    body?: TRequestBody
    params?: TRequestParams
    /** token for request cancel */
    cancelToken?: CancelToken
}

/** {@link THook} extends this */
export type TBasicHook = {
    config: TRequestConfig
}

/** used in core for hooks executing */
export type THook = TBasicHook & {
    name: "onRequest"
    data: TRequestConfig
} | TBasicHook & {
    name: "onResponse"
    data: TResponse
} | TBasicHook & {
    name: "onDownload" | "onUploadProgress" | "onUploaded"
    data: ProgressEvent<EventTarget>
} | TBasicHook & {
    name: "onError"
    data: TError
}

/** execute functions on XHR lifecycle */
export interface IHooks {
    /** request sended */
    onRequest?: (config: TRequestConfig) => void
    /** response came */
    onResponse?: (response: TResponse) => void
    /** downloading data from server */
    onDownload?: (e: ProgressEvent<EventTarget>) => void
    /** uploading to server */
    onUploadProgress?: (e: ProgressEvent<EventTarget>) => void
    /** uploaded to server */
    onUploaded?: (e: ProgressEvent<EventTarget>) => void
    /** request (like 404) / network (like CORS) / timeout error */
    onError?: (err: TError) => void
}