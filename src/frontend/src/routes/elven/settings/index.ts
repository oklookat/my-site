import type { RequestHandlerOutput } from "@sveltejs/kit";
import type { RequestEvent } from "@sveltejs/kit/types/private";
//

export async function get(event: RequestEvent): Promise<RequestHandlerOutput> {
    return {
        body: {
            user: event.locals.user
        }
    }
}