import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//
import NetworkFile from "$lib_elven/network/network_file";
import ToolsFiles from "$lib_elven/tools/files";
import { searchParamsToObject } from "$lib_elven/tools";


/** get files */
export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    const defaultParams = ToolsFiles.getDefaultParams()
    let requestParams = searchParamsToObject(event.url.searchParams)
    requestParams = Object.assign(defaultParams, requestParams)

    const networkFile = new NetworkFile(event.locals.user.token)
    const resp = await networkFile.getAll(requestParams)
    let items = {}

    if (resp.ok) {
        items = await resp.json()
    }

    return {
        // @ts-ignore
        body: {
            items: items,
            params: requestParams
        }
    };
}