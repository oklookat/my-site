import Duck from './Duck'
import type { TArticle, TData, TParams } from '@/types/ArticleTypes'

class ArticleAdapter {

    public static async getAll(params: TParams): Promise<TData> {
        try {
            const response = await Duck.GET({ url: 'articles', params: params })
            return Promise.resolve(response.body as TData)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async get(id: string): Promise<TArticle> {
        try {
            const response = await Duck.GET({ url: `articles/${id}` })
            return Promise.resolve(response.body as TArticle)
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

    public static async create(article: TArticle): Promise<TArticle> {
        try {
            const response = await Duck.POST({ url: `articles`, body: article })
            return Promise.resolve(response.body as TArticle)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async update(article: TArticle): Promise<TArticle> {
        try {
            const response = await Duck.PATCH({ url: `articles/${article.id}`, body: article })
            return Promise.resolve(response.body as TArticle)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async publish(id: string): Promise<TArticle> {
        try {
            const response = await Duck.PATCH({ url: `articles/${id}`, body: { is_published: true } })
            return Promise.resolve(response.body as TArticle)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async makeDraft(id: string): Promise<TArticle> {
        try {
            const response = await Duck.PATCH({ url: `articles/${id}`, body: { is_published: false } })
            return Promise.resolve(response.body as TArticle)
        } catch (err) {
            return Promise.reject(err)
        }
    }
}

export default ArticleAdapter