import {Axios} from './Axios'
import type { IArticle, IArticlesData } from '@/types/ArticleTypes'

class ArticleAdapter {

    public static async getAll(page = 1, show = 'published', sortBy = 'updated', start = 'newest'): Promise<IArticlesData> {
        const config = {
            params:
                {
                    page: page, show: show, by: sortBy, start: start
                }
        }
        try {
            const response = await Axios.get('articles', config)
            return Promise.resolve(response.data as IArticlesData)
        } catch(err) {
            return Promise.reject(err)
        }
    }

    public static async get(id: string): Promise<IArticle> {
        try {
            const response = await Axios.get(`articles/${id}`)
            return Promise.resolve(response.data as IArticle)
        } catch (err) {
            return Promise.reject(err);
        }
    }

    public static async delete(id: string) {
        try {
            await Axios.delete(`articles/${id}`)
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err);
        }
    }

    public static async create(article) {
        try {
            const response = await Axios.post('articles', article)
            return Promise.resolve(response.data as IArticle)
        } catch(err) {
            return Promise.reject(err)
        }
    }

    public static async update(article) {
        try {
            const response = await Axios.put(`articles/${article.id}`, article)
            return Promise.resolve(response.data as IArticle)
        } catch (err){
            return Promise.reject(err)
        }
    }

    public static async publish(article) {
        try {
            const response = await Axios.put(`articles/${article.id}`, {is_published: true})
            return Promise.resolve(response.data as IArticle)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async makeDraft(article) {
        try {
            const response = await Axios.put(`articles/${article.id}`, {is_published: false})
            return Promise.resolve(response.data as IArticle)
        } catch (err) {
            return Promise.reject(err)
        }
    }
}

export default ArticleAdapter