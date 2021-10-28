import Duck from './Duck'
import type { IArticle, IArticlesData } from '@/types/ArticleTypes'

class ArticleAdapter {

    public static async getAll(page = 1, show = 'published', sortBy = 'updated', start = 'newest'): Promise<IArticlesData> {
        const params = {page: page, show: show, by: sortBy, start: start}
        try {
            const response = await Duck.GET({url: 'articles', params})
            return Promise.resolve(response.body as IArticlesData)
        } catch(err) {
            return Promise.reject(err)
        }
    }

    public static async get(id: string): Promise<IArticle> {
        try {
            const response = await Duck.GET({url: `articles/${id}`})
            return Promise.resolve(response.body as IArticle)
        } catch (err) {
            return Promise.reject(err);
        }
    }

    public static async delete(id: string) {
        try {
            await Duck.DELETE({url: `articles/${id}`})
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err);
        }
    }

    public static async create(article) {
        try {
            const response = await Duck.POST({url: `articles`, body: article})
            return Promise.resolve(response.body as IArticle)
        } catch(err) {
            return Promise.reject(err)
        }
    }

    public static async update(article) {
        try {
            const response = await Duck.PUT({url: `articles/${article.id}`, body: article})
            return Promise.resolve(response.body as IArticle)
        } catch (err){
            return Promise.reject(err)
        }
    }

    public static async publish(article) {
        try {
            const response = await Duck.PUT({url: `articles/${article.id}`, body: {is_published: true}})
            return Promise.resolve(response.body as IArticle)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async makeDraft(article) {
        try {
            const response = await Duck.PUT({url: `articles/${article.id}`, body: {is_published: false}})
            return Promise.resolve(response.body as IArticle)
        } catch (err) {
            return Promise.reject(err)
        }
    }
}

export default ArticleAdapter