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
    type: "timeout" | "network" | "response"
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
export interface IHooks {
    // request hooks
    onRequest?: (config: TRequestConfig) => void
    onRequestError?: (err: TError) => void
    // response hooks
    onResponse?: (response: TResponse) => void
    onResponseError?: (err: TError) => void
    // upload hooks
    onUploadProgress?: () => void
    onUploaded?: () => void
    onUploadError?: (err: TError) => void
}