// basic
export type TRequestMethod = "GET" | "POST" | "PUT" | "DELETE" | "HEAD" | "OPTIONS" | "PATCH"
export type TRequestBody = string | number | object | Blob | BufferSource | FormData | URLSearchParams | ReadableStream

export type THeaders = {
    [name: string]: string | number
}

export type TRequestParams = {
    [name: string]: string | number
}

export type TResponse = { 
    body: any
    statusCode: number
}

export type TError = TResponse & {
    type: "timeout" | "network" | "request" | "response"
    details?: ProgressEvent<EventTarget>
}

export interface PromiseWithCancel extends Promise<TResponse> {
    cancel: (xhr: XMLHttpRequest) => void
}

// config
export type TBasicConfig = {
    withCredentials?: boolean
    headers?: THeaders
    hooks?: IHooks
}

export type TGlobalConfig = TBasicConfig & {
    timeout?: number
    baseURL?: string
}

export type TRequestConfig = TBasicConfig & {
    url: string
    body?: TRequestBody
    params?: TRequestParams
}

// hooks
export type TBasicHook = {
    config: TRequestConfig
}
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

export interface IHooks {
    // request hooks
    onRequest?: (config: TRequestConfig) => void
    // response hooks
    onResponse?: (response: TResponse) => void
    // download from server hooks
    onDownload?: (e: ProgressEvent<EventTarget>) => void
    // upload to server hooks
    onUploadProgress?: (e: ProgressEvent<EventTarget>) => void
    onUploaded?: (e: ProgressEvent<EventTarget>) => void
    // error
    onError?: (err: TError) => void
}