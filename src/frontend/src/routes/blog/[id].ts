import NetworkArticle from "$lib_elven/network/network_article";
import type { Article } from "$lib_elven/types/articles";
import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";

export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    let resp: Response
    let article: Article | null = null
    let statusCode = 200
    try {
        const networkArticle = new NetworkArticle('')
        resp = await networkArticle.get(event.params.id)
        if (resp.ok) {
            article = await resp.json()
        } else {
            statusCode = resp.status
        }
    } catch (err) {
        statusCode = 500
    }

    return {
        status: statusCode,
        body: { article: article }
    };
}