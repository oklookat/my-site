import { stringToNormal } from "$lib_elven/tools";
import type { Article } from "$lib_elven/types/articles";
import { By as FilesBy, Start, type File, type Params as FileParamsT } from "$lib_elven/types/files";
import { type Params as ArticleParamsT, By as ArticlesBy } from "$lib_elven/types/articles";

type Param_Type<T extends any> =
    T extends File ? FileParamsT :
    T extends Article ? ArticleParamsT : never;

type ForWhom = "file" | "article"

/** manage request params (file, article) */
export class Params<T> {

    private forWhom: ForWhom
    private self: Param_Type<T>
    private searchparams: URLSearchParams

    /** create new request params */
    constructor(forWhom: ForWhom, init?: URLSearchParams) {
        if (forWhom !== 'article' && forWhom !== 'file') {
            throw Error("wrong 'forWhom' param")
        }

        this.forWhom = forWhom
        this.self = this.getDefault()
        this.import(init)
    }

    /** import values from searchparams */
    public import(params: URLSearchParams) {
        // @ts-ignore
        this.searchparams = new URLSearchParams(this.self)
        
        if (!params || !(params instanceof URLSearchParams)) {
            return
        }

        this.searchparams = params

        this.searchparams.forEach((value, key) => {
            const normalVal = stringToNormal(value)
            // @ts-ignore
            this.setParam(key, normalVal)
        })
    }

    /** set request param */
    public setParam(name: keyof Param_Type<T>, val: any) {
        // @ts-ignore
        name = String(name)
        // @ts-ignore
        name = name.toLowerCase()

        const isExists = name in this.getDefault()
        if (!isExists) {
            return
        }

        const normalized = stringToNormal(val)
        if (normalized === undefined || normalized === null || normalized === '') {
            // @ts-ignore
            this.searchparams.delete(name)
            delete this.self[name]
            return
        }

        // @ts-ignore
        if (name === 'page') {
            if (val < 1 || isNaN(Number(val))) {
                // @ts-ignore
                val = 1
            }
        }

        // @ts-ignore
        this.searchparams.set(name, val)

        // @ts-ignore
        this.self[name] = val
    }

    /** get request param */
    public getParam(name: keyof Param_Type<T>): any {
        return this.self[name]
    }

    /** get request params searchparams copy */
    public toSearchparams(): URLSearchParams {
        return new URLSearchParams(this.searchparams)
    }

    /** get request params copy */
    public toObject(): Param_Type<T> {
        return { ...this.self }
    }

    /** get default request params */
    private getDefault(): Param_Type<T> {
        if (this.forWhom === 'article') {
            return {
                page: 1,
                drafts: false,
                newest: true,
                preview: true,
                by: ArticlesBy.published,
                title: undefined
            } as Param_Type<T>
        } else if (this.forWhom === 'file') {
            return {
                page: 1,
                start: Start.newest,
                by: FilesBy.created,
                extensions: undefined,
                filename: undefined
            } as Param_Type<T>
        }
    }
}