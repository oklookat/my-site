import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//
import NetworkArticle from "$lib/network/network_article";
import { getDefaultParams, type Params } from "$lib/types/articles";
import Utils from "$lib/tools";


export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    const requestParams = getDefaultParams()
    //
    const paramsObj = Utils.searchParamsToObject(event.url.searchParams)
    Object.assign(requestParams, paramsObj)

    const networkArticle = new NetworkArticle(event.locals.user.token)
    const items = await networkArticle.getAll(requestParams)
    return {
        body: { 
            items: items, 
            params: requestParams 
        }
    };
}