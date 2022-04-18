import NetworkCategories from "$lib_elven/network/network_category";
import type { Items } from "$lib_elven/types";
import type { Category } from "$lib_elven/types/articles/categories";
import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//

export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    const networkCategories = new NetworkCategories(event.locals.user.token)
    let items: Items<Category>
    const resp = await networkCategories.getAll()
    if(resp.status === 200) {
        items = await resp.json()
    }
    return {
        body: {
            items: items
        }
    }

}