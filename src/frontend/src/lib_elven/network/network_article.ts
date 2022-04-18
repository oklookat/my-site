import Fetchd from '$lib_elven/network'
import Utils from '$lib_elven/tools'
import type { Article, Params } from '$lib_elven/types/articles'

/** Use with SSR by passing token / or in components by passing empty token.
 * 
 * Static methods = not for SSR, use only on components side
 */
export default class NetworkArticle {

    private static prefix = "article/articles"
    private headers: Headers

    constructor(token: string) {
        const headers = new Headers()
        Utils.addTokenToHeaders(token, headers)
        this.headers = headers
    }

    public async getAll(params: Params): Promise<Response> {

        try {
            const response = await Fetchd.send({
                method: "GET",
                url: NetworkArticle.prefix, params: params, headers: this.headers
            })
            return response
        } catch (err) {
            throw err
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
            throw err
        }
    }

    public static async delete(id: string): Promise<Response> {
        try {
            const resp = await Fetchd.send({
                method: "DELETE",
                url: `${this.prefix}/${id}`})
            return resp
        } catch (err) {
            throw err
        }
    }

    public static async create(article: Article): Promise<Response> {
        NetworkArticle.beforeCRUD(article)
        try {
            const response = await Fetchd.send({
                method: "POST",
                url: `${this.prefix}`, body: article})
            return response
        } catch (err) {
            throw err
        }
    }

    public static async update(article: Article): Promise<Response> {
        NetworkArticle.beforeCRUD(article)
        try {
            const response = await Fetchd.send({
                method: "PATCH", url: `${this.prefix}/${article.id}`,
                body: article})
            return response
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