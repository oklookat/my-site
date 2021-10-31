// from: https://github.com/axios/axios/blob/76f09afc03fbcf392d31ce88448246bcd4f91f8c/lib/utils.js

export class Validator {

    public toString = Object.prototype.toString

    // basic
    public static isUndefined(val: any): boolean {
        return typeof val === 'undefined'
    }

    public static isObject(val: any): boolean {
        return val !== null && typeof val === 'object'
    }

    public static isFunction(val: any): boolean {
        return val instanceof Function
    }

    public static isArray(val: any): boolean {
        return this.toString.call(val) === '[object Array]'
    }

    public static isString(val): boolean {
        return typeof val === 'string'
    }

    // main
    public static isFormData(val: any): boolean {
        return val instanceof FormData
    }

    public static isArrayBuffer(val: any): boolean {
        return toString.call(val) === '[object ArrayBuffer]'
    }

    public static isBuffer(val: any): boolean {
        return val !== null && !this.isUndefined(val) && val.constructor !== null && !this.isUndefined(val.constructor)
            && typeof val.constructor.isBuffer === 'function' && val.constructor.isBuffer(val)
    }

    public static isStream(val: any): boolean {
        return this.isObject(val) && this.isFunction(val.pipe)
    }

    public static isFile(val: any): boolean {
        return toString.call(val) === '[object File]'
    }

    public static isBlob(val: any): boolean {
        return val instanceof Blob
    }

    public static isURLSearchParams(val: any): boolean {
        return typeof URLSearchParams !== 'undefined' && val instanceof URLSearchParams
    }

    public static isArrayBufferView(val: any): boolean {
        let result: any
        if (ArrayBuffer.isView) {
            result = ArrayBuffer.isView(val)
        } else {
            result = (val) && (val.buffer) && (val.buffer instanceof ArrayBuffer)
        }
        return result
    }
}