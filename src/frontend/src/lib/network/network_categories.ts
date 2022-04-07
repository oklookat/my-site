import Fetchd from '$lib/network'
import { StorageAuth } from '$lib/tools/storage'
import type { Data } from '$lib/types'
import type { Category } from '$lib/types/articles/categories'

export default class NetworkCategories {

    private static prefix = "article/categories"
    private headers: Headers

    constructor(token: string) {
        const headers = new Headers()
        StorageAuth.addTokenToHeaders(headers, token)
        this.headers = headers
    }

    public async getAll() {
        try {
            const response = await Fetchd.send({ method: "GET", url: NetworkCategories.prefix, headers: this.headers})
            const jsond = await response.json()
            return jsond as Data<Category>
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public async create(cat: Category) {
        try {
            const response = await Fetchd.send({ 
                method: "POST", url: NetworkCategories.prefix, body: cat,
                headers: this.headers
            })
            const jsond = await response.json()
            return jsond as Category
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public async rename(cat: Category) {
        try {
            const response = await Fetchd.send({ 
                method: "PATCH", url: `${NetworkCategories.prefix}/${cat.id}`, 
                body: { name: cat.name },
                headers: this.headers
            })
            const jsond = await response.json()
            return jsond as Category
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public async delete(id: string) {
        try {
            const response = await Fetchd.send({ 
                method: "DELETE", url: `${NetworkCategories.prefix}/${id}`,
                headers: this.headers
            })
            return
        } catch (err) {
            throw err
        }
    }

}