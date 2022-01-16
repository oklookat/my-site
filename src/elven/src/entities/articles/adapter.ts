import Duck from '@/duck'
import type { Data } from '@/types'
import type { Article, Params } from '@/entities/articles/types'

export default class ArticleAdapter {

    public static async getAll(params: Params): Promise<Data<Article>> {
        try {
            const response = await Duck.GET({ url: 'articles', params: params })
            return Promise.resolve(response.body as Data<Article>)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async get(id: string): Promise<Article> {
        try {
            const response = await Duck.GET({ url: `articles/${id}` })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err);
        }
    }

    public static async delete(id: string): Promise<void> {
        try {
            await Duck.DELETE({ url: `articles/${id}` })
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err);
        }
    }

    public static async create(article: Article): Promise<Article> {
        try {
            const response = await Duck.POST({ url: `articles`, body: article })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async update(article: Article): Promise<Article> {
        console.log(article.content)
        try {
            const response = await Duck.PATCH({ url: `articles/${article.id}`, body: article })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async publish(id: string): Promise<Article> {
        try {
            const response = await Duck.PATCH({ url: `articles/${id}`, body: { is_published: true } })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async unpublish(id: string): Promise<Article> {
        try {
            const response = await Duck.PATCH({ url: `articles/${id}`, body: { is_published: false } })
            return Promise.resolve(response.body as Article)
        } catch (err) {
            return Promise.reject(err)
        }
    }

}