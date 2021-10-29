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

// config
export type TBasicConfig = {
    withCredentials?: boolean
    headers?: THeaders
    hooks?: IHooks
    // base auth
    baseAuth?: boolean
    baseAuthUser?: string
    baseAuthPassword?: string
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
export type THook = {
    name: "onRequest" | "onResponse" | "onDownload" | "onUploadProgress" | "onUploaded" | "onError"
    data: string
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