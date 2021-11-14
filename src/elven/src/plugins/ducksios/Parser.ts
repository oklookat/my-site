// some from: https://github.com/axios/axios/blob/76f09afc03fbcf392d31ce88448246bcd4f91f8c/dist/axios.js#L1267
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

import Service from "./service";
import type { TGlobalConfig, TRequestConfig } from "./types";
import { Validator } from "./validator";

/** parse body before request and after response */
export default class Parser {

    /** set headers if not set. If it's json auto stringify before request */
    public static requestBody(body: any, rc: TRequestConfig, gc: TGlobalConfig): { r: TRequestConfig, b: any } {
        const simpleData = Validator.isFormData(body) ||
            Validator.isArrayBuffer(body) ||
            Validator.isBuffer(body) ||
            Validator.isStream(body) ||
            Validator.isFile(body) ||
            Validator.isBlob(body)
        if (simpleData) {
            return { r: rc, b: body }
        }
        if (Validator.isArrayBufferView(body)) {
            return { r: rc, b: body.buffer }
        }
        if (Validator.isURLSearchParams(body)) {
            rc = Service.setContentTypeIfUnset('application/x-www-form-urlencoded; charset=utf-8', gc, rc)
            return { r: rc, b: body.toString() }
        }
        // if json convertable
        if (Validator.isObject(body)) {
            rc = Service.setContentTypeIfUnset('application/json', gc, rc)
            return { r: rc, b: Service.stringifySafely(body) }
        }
        return { r: rc, b: body }
    }

    /** auto parse json after response */ 
    public static responseBody(body: any): any {
        if (!body) {
            return body
        }
        try {
            return JSON.parse(body);
        } catch (e) {
        }
        return body
    }

}