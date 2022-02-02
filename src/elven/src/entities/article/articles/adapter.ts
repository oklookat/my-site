import Duck from '@/duck'
import type { Data } from '@/types'
import type { Article, Params } from './types'

export default class ArticleAdapter {

    private static prefix = "article/articles"

    public static async getAll(params: Params): Promise<Data<Article>> {
        try {
            const response = await Duck.GET({ url: this.prefix, params: params })
            return Promise.resolve(response.body as Data<Article>)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async get(id: string): Promise<Article> {
        try {
            const response = await Duck.GET({ url: `${this.prefix}/${id}` })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err);
        }
    }

    public static async delete(id: string): Promise<void> {
        try {
            await Duck.DELETE({ url: `${this.prefix}/${id}` })
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err);
        }
    }

    public static async create(article: Article): Promise<Article> {
        try {
            const response = await Duck.POST({ url: `${this.prefix}`, body: article })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async update(article: Article): Promise<Article> {
        try {
            const response = await Duck.PATCH({ url: `${this.prefix}/${article.id}`, body: article })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async publish(id: string): Promise<Article> {
        try {
            const response = await Duck.PATCH({ url: `${this.prefix}/${id}`, body: { is_published: true } })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async unpublish(id: string): Promise<Article> {
        try {
            const response = await Duck.PATCH({ url: `${this.prefix}/${id}`, body: { is_published: false } })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err)
        }
    }

}