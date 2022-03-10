import Duck from '@/network'
import type { Data } from '@/types'
import type { Category } from './types'

export default class CategoriesNetwork {

    private static prefix = "article/categories"

    public static async getAll() {
        try {
            const response = await Duck.GET({ url: this.prefix })
            return Promise.resolve(response.body as Data<Category>)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async create(cat: Category) {
        try {
            const response = await Duck.POST({ url: `${this.prefix}`, body: cat })
            return Promise.resolve(response.body as Category)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async rename(cat: Category) {
        try {
            const response = await Duck.PATCH({ url: `${this.prefix}/${cat.id}`, body: { name: cat.name } })
            return Promise.resolve(response.body as Category)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async delete(id: string) {
        try {
            await Duck.DELETE({ url: `${this.prefix}/${id}` })
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err);
        }
    }

}