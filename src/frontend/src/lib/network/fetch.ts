import {FormData} from "formdata-node"
//
export type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
export type Body = string | number | object | Blob | BufferSource | FormData | URLSearchParams | ReadableStream
export interface Stringer {
    toString: () => string
}
export type RequestConfig = {
    method: RequestMethod,
    url: string,
    body?: Body
    params?: Record<string | number, Stringer>
    headers?: Headers
}

export default class FetchDriver {
    
    private baseURL: string

    constructor(baseURL: string) {
        this.baseURL = baseURL
    }

    public async send(config: RequestConfig): Promise<Response> {
        let endpoint = ""

        if(this.baseURL) {
            endpoint = `${this.baseURL}/${config.url}`
        } else {
            endpoint = config.url
        }
        const url = new URL(endpoint)
        if(config.params) {
            for(const key in config.params) {
                const val = config.params[key]
                if(!val) {
                    continue
                }
                // @ts-ignore
                url.searchParams.append(key, val)
            }
        }

        // TODO: if formdata, auto-set content-type/boundary headers not work. Need fix
        var headers = new Headers();
        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        const isFormData = config.body && config.body instanceof FormData
        if(isFormData) {
            headers.delete('Content-Type');
        }

        if(config.headers) {
            config.headers.forEach((val, key) => {
                headers.append(key, val)
            })
        }

        const fetchConfig: RequestInit = {
            method: config.method,
            credentials: 'include',
            headers: headers,
            body: JSON.stringify(config.body || {})
        }

        if (config.method == "GET") {
            delete fetchConfig.body
        }
        
        try {
            const result = await fetch(url.toString(), fetchConfig)
            return result
        } catch(err) {
            return Promise.reject(err)
        }
    }
}