import NetworkCategories from "$lib/network/network_categories";
import type { Category } from "$lib/types/articles/categories";
import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//

export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    const networkCategories = new NetworkCategories(event.locals.user.token)
    const items = await networkCategories.getAll()
    return {
        body: {items: items}
    }

}

export async function post(event: RequestEvent): Promise<RequestHandlerOutput> {
    const cat: Category = await event.request.json()
    const networkCategories = new NetworkCategories(event.locals.user.token)
    await networkCategories.create(cat)
    return {}
}

export async function del(event: RequestEvent): Promise<RequestHandlerOutput> {
    const id = await event.request.text()
    const networkCategories = new NetworkCategories(event.locals.user.token)
    await networkCategories.delete(id)
    return {}
}

export async function patch(event: RequestEvent): Promise<RequestHandlerOutput> {
    const cat: Category = await event.request.json()
    const networkCategories = new NetworkCategories(event.locals.user.token)
    await networkCategories.rename(cat)
    return {}
}