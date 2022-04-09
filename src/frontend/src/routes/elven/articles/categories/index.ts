import NetworkCategories from "$lib/network/network_category";
import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//

export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    const networkCategories = new NetworkCategories(event.locals.user.token)
    const items = await networkCategories.getAll()
    return {
        body: {
            items: items
        }
    }

}