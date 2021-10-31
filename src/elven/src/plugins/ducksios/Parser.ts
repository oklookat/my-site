// some from: https://github.com/axios/axios/blob/76f09afc03fbcf392d31ce88448246bcd4f91f8c/dist/axios.js#L1267

import Service from "./Service";
import type { TGlobalConfig, TRequestConfig } from "./Types";
import { Validator } from "./Validator";

export default class Parser {

    // set headers if not set, and if json, auto stringify it before request
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

    // auto parse json after response
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