import NetworkFile from "$lib/network/network_files";
import Utils from "$lib/tools";
import { getDefaultParams } from "$lib/types/files";
import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//


/** get files */
export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    const requestParams = getDefaultParams()
    //
    const paramsObj = Utils.searchParamsToObject(event.url.searchParams)
    Object.assign(requestParams, paramsObj)

    const networkFile = new NetworkFile(event.locals.user.token)
    const items = await networkFile.getAll(requestParams)
    return {
        body: { 
            items: items, 
            params: requestParams 
        }
    };
}

/** upload file */
export async function post(event: RequestEvent): Promise<RequestHandlerOutput> {
    const file = await event.request.json()
    const networkFile = new NetworkFile(event.locals.user.token)
    await networkFile.upload(file)
    return {}
}