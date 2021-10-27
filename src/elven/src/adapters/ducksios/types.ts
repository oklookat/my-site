// basic types
export type TRequestMethod = "GET" | "POST" | "PUT" | "DELETE" | "HEAD" | "OPTIONS" | "PATCH"

export type THeader = {
    name: string
    value: string
}

export type TRequestParam = {
    name: string
    value: string | number
}

export type TRequestBody = string | number | object | Blob | BufferSource | FormData | URLSearchParams | ReadableStream

export type TResponse = { 
    body: any
    statusCode: number
}

export interface IError extends TResponse {
    type: "timeout" | "network" | "response"
}

// hooks
export interface IHooks {
    onRequest?: () => null
    onRequestError?: (err: IError) => null

    onResponse?: () => null
    onResponseError?: (err: IError) => null

    onUploadProgress?: () => null
    onUploaded?: () => null
    onUploadError?: (err: IError) => null
}

// config
export type TBasicConfig = {
    withCredentials?: boolean
    headers?: Array<THeader>
    hooks?: IHooks
    //
    baseAuth?: boolean
    baseAuthUser?: string
    baseAuthPassword?: string
}

export interface IGlobalConfig extends TBasicConfig {
    timeout?: number
    baseURL?: string
}

export interface IRequestConfig extends TBasicConfig {
    url: string
    body?: TRequestBody
    params?: Array<TRequestParam>
}