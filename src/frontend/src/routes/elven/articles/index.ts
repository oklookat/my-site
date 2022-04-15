import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//
import NetworkArticle from "$lib_elven/network/network_article";
import Utils from "$lib_elven/tools";
import ToolsArticles from "$lib_elven/tools/articles";


export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    const requestParams = ToolsArticles.getDefaultParams()
    //
    const paramsObj = Utils.searchParamsToObject(event.url.searchParams)
    Object.assign(requestParams, paramsObj)

    const networkArticle = new NetworkArticle(event.locals.user.token)
    const resp = await networkArticle.getAll(requestParams)
    let items = {}
    if (resp.status === 200) {
        items = await resp.json()
    }
    return {
        body: {
            items: items,
            params: requestParams
        }
    };
}