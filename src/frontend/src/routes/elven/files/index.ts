import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//
import { getDefaultParams } from "$lib/types/files";
import NetworkFile from "$lib/network/network_file";
import Utils from "$lib/tools";


/** get files */
export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    const requestParams = getDefaultParams()
    //
    const paramsObj = Utils.searchParamsToObject(event.url.searchParams)
    Object.assign(requestParams, paramsObj)

    const networkFile = new NetworkFile(event.locals.user.token)
    const items = await networkFile.getAll(requestParams)
    return {
        // @ts-ignore
        body: { 
            items: items, 
            params: requestParams 
        }
    };
}