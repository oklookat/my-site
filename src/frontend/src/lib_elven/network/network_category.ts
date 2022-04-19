import Fetchd from '$lib_elven/network'
import Utils from '$lib_elven/tools'
import type { Items } from '$lib_elven/types'
import type { Category } from '$lib_elven/types/articles/categories'

/** Use with SSR by passing token / or in components by passing empty token.
 * 
 * Static methods = not for SSR, use only on components side
 */
export default class NetworkCategory {

    private static prefix = "article/categories"
    private headers: Headers

    constructor(token: string) {
        const headers = new Headers()
        Utils.addTokenToHeaders(token, headers)
        this.headers = headers
    }

    public async getAll(): Promise<Response> {
        try {
            const response = await Fetchd.send({ method: "GET", url: NetworkCategory.prefix, headers: this.headers })
            return response
        } catch (err) {
            throw err
        }
    }

    /** get category by name */
    public async get(name: string): Promise<Response> {
        try {
            const response = await Fetchd.send({ method: "GET", url: NetworkCategory.prefix + `/${name}`, headers: this.headers })
            return response
        } catch (err) {
            throw err
        }
    }

    public static async create(cat: Category): Promise<Response> {
        try {
            const response = await Fetchd.send({
                method: "POST", url: this.prefix, body: cat
            })
            return response
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async rename(cat: Category): Promise<Response> {
        try {
            const response = await Fetchd.send({
                method: "PATCH", url: `${this.prefix}/${cat.id}`,
                body: { name: cat.name }
            })
            return response
        } catch (err) {
            throw err
        }
    }

    public static async delete(id: string) {
        try {
            const response = await Fetchd.send({
                method: "DELETE", url: `${this.prefix}/${id}`
            })
            return response
        } catch (err) {
            throw err
        }
    }

}