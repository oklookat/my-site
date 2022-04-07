import Fetchd from '$lib/network'
import { StorageAuth } from '$lib/tools/storage'
import type { Data } from '$lib/types'
import type { Article, Params } from '$lib/types/articles'

export default class NetworkArticle {

    private static prefix = "article/articles"
    private headers: Headers

    constructor(token: string) {
        const headers = new Headers()
        StorageAuth.addTokenToHeaders(headers, token)
        this.headers = headers
    }

    public async getAll(params: Params): Promise<Data<Article>> {

        try {
            const response = await Fetchd.send({
                method: "GET",
                url: NetworkArticle.prefix, params: params, headers: this.headers
            })
            const jsond = await response.json()
            return jsond as Data<Article>
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public async get(id: string): Promise<Article> {
        try {
            const response = await Fetchd.send({
                method: "GET",
                url: `${NetworkArticle.prefix}/${id}`, headers: this.headers
            })
            const jsond = await response.json()
            return jsond as Article
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public async delete(id: string): Promise<void> {
        try {
            await Fetchd.send({
                method: "DELETE",
                url: `${NetworkArticle.prefix}/${id}`, headers: this.headers
            })
            return
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public async create(article: Article): Promise<Article> {
        NetworkArticle.beforeCRUD(article)
        try {
            const response = await Fetchd.send({
                method: "POST",
                url: `${NetworkArticle.prefix}`, body: article, headers: this.headers
            })
            const jsond = await response.json()
            return jsond as Article
        } catch (err) {
            throw err
        }
    }

    public async update(article: Article): Promise<Article> {
        NetworkArticle.beforeCRUD(article)
        try {
            const response = await Fetchd.send({
                // PATCH OR POST?
                method: "PATCH", url: `${NetworkArticle.prefix}/${article.id}`,
                body: article, headers: this.headers
            })
            const jsond = await response.json()
            return jsond as Article
        } catch (err) {
            throw err
        }
    }

    private static beforeCRUD(article: Article) {
        if (!article.category_id) {
            article.category_id = "nope"
        }
    }
}