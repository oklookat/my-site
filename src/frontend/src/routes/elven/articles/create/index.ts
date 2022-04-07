import NetworkArticle from "$lib/network/network_article";
import type { Article } from "$lib/types/articles";
import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";

export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    /** creating / editing this article */
    let article: Article = {
        title: '',
        content: ''
    }

    const params = event.url.searchParams
    const isEditMode = params.has("id")
    if (!isEditMode) {
        return {
            body: { article: article }
        }
    }
    try {
        const networkArticle = new NetworkArticle(event.locals.user.token)
        const articled = await networkArticle.get(params.get("id"))
        article = articled
    } catch (err) { }
    return {
        body: { article: article }
    }
}

export async function post(event: RequestEvent): Promise<RequestHandlerOutput> {
    const article = await event.request.json()
    const networkArticle = new NetworkArticle(event.locals.user.token)
    const newArticle = await networkArticle.create(article)
    return {
        body: { newArticle }
    }
}

export async function patch(event: RequestEvent): Promise<RequestHandlerOutput> {
    const article = await event.request.json()
    const networkArticle = new NetworkArticle(event.locals.user.token)
    const updatedArticle = await networkArticle.update(article)
    return {
        body: { updatedArticle }
    }
}

export async function del(event: RequestEvent): Promise<RequestHandlerOutput> {
    const articleID = await event.request.text()
    const networkArticle = new NetworkArticle(event.locals.user.token)
    await networkArticle.delete(articleID)
    return {}
}