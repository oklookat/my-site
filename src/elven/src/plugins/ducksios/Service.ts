// some from: https://github.com/axios/axios/blob/76f09afc03fbcf392d31ce88448246bcd4f91f8c/lib/utils.js
// MIT license stuff (axios):
/**
Copyright (c) 2014-present Matt Zabriskie
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
 */

import type { TGlobalConfig, TRequestConfig, THeaders, TRequestParams } from "./types";
import { Validator } from "./validator";


export default class Service {

    /**
     * set base URL
     * @param url url like: /user
     * @param base url like: example.com
     * @returns url like: example.com/user
     */
    public static setBaseURL(url: string, base?: string): string {
        if (!base) {
            return url
        }
        return Service.replaceDoubleSlashes(`${base}/${url}`)
    }

    /**
     * replace double slashes
     * @param s string like: hello//world////
     * @returns string like: hello/world
     */
    public static replaceDoubleSlashes(s: string): string {
        return s.replace(/(?<!:)\/\/+/g, '/')
    }

    /** trim excess whitespace off the beginning and end of a string */
    public static trim(str: string): string {
        str = str.trim()
        str = str.replace(/^\s+|\s+$/g, '')
        return str
    }

    /** set request params to url using {@link URL} object */
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

    /** set request headers to XHR based on settings */
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

    /** check content type in local and global config. If content type not in local and global, set header in local config */
    public static setContentTypeIfUnset(value: string | number, gc: TGlobalConfig, rc: TRequestConfig): TRequestConfig {
        const notInGlobal = gc.headers && !('Content-Type' in gc.headers)
        const notInLocal = rc.headers && !('Content-Type' in rc.headers)
        if (notInLocal && notInGlobal) {
            rc.headers['Content-Type'] = value
        }
        return rc
    }

    /** if rawValue string - check is valid json, trim and return. Otherwise - stringify and return. If error while validation or stringify - throws error */
    public static stringifySafely(rawValue: any): any {
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