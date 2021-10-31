// some from: https://github.com/axios/axios/blob/76f09afc03fbcf392d31ce88448246bcd4f91f8c/lib/utils.js

import type { TGlobalConfig, TRequestConfig, THeaders, TRequestParams } from "./Types";
import { Validator } from "./Validator";

export default class Service {

    public static setBaseURL(url: string, base?: string): string {
        if (!base) {
            return url
        }
        return Service.replaceDoubleSlashes(`${base}/${url}`)
    }

    public static replaceDoubleSlashes(s: string): string {
        return s.replace(/(?<!:)\/\/+/g, '/')
    }

    // trim excess whitespace off the beginning and end of a string
    public static trim(str: string): string {
        str = str.trim()
        str = str.replace(/^\s+|\s+$/g, '')
        return str
    }

    // set request params to url
    public static setRequestParams(url: string, params?: TRequestParams): string {
        if (!params) {
            return url
        }
        let urlObj = new URL(url)
        for (const param in params) {
            urlObj.searchParams.set(param, params[param].toString())
        }
        return urlObj.toString()
    }

    // set request headers to xhr based on configs
    public static setRequestHeaders(xhr: XMLHttpRequest, rc: TRequestConfig, gc: TGlobalConfig): XMLHttpRequest {
        const set = (headers: THeaders) => {
            for (const header in headers) {
                xhr.setRequestHeader(header, headers[header].toString())
            }
        }
        let headers = gc.headers
        if (headers) {
            set(gc.headers)
        }
        headers = rc.headers
        if (headers) {
            set(rc.headers)
        }
        return xhr
    }

    // check content type in local and global config. If content type not in local and global, set header in local config
    public static setContentTypeIfUnset(value: string | number, gc: TGlobalConfig, rc: TRequestConfig): TRequestConfig {
        const notInGlobal = gc.headers && !('Content-Type' in gc.headers)
        const notInLocal = rc.headers && !('Content-Type' in rc.headers)
        if (notInLocal && notInGlobal) {
            rc.headers['Content-Type'] = value
        }
        return rc
    }

    // if rawValue string - check is valid json; trim and return. Otherwise - stringify and return. Or error while validation or stringify - throws error.
    public static stringifySafely(rawValue: any) {
        if (Validator.isString(rawValue)) {
            try {
                JSON.parse(rawValue)
                return this.trim(rawValue)
            } catch (err) {
                if (err.name !== 'SyntaxError') {
                    throw err
                }
            }
        }
        try {
            rawValue = JSON.stringify(rawValue)
            return rawValue
        } catch (err) {
            throw err
        }
    }
}