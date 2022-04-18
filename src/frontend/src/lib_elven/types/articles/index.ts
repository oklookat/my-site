import type { Page } from "$lib_elven/types"

/** request param */
export enum By {
    created = 'created',
    published = 'published',
    updated = 'updated'
}

/** request params */
export type Params = {
    page?: Page

    published?: boolean

    newest?: boolean

    preview?: boolean

    /** show only without categories */
    without_category?: boolean

    by?: By

    /** search by category name */
    category_name?: string

    /** search by title */
    title?: string
}

/** article */
export type Article = {
    id?: string
    user_id?: string
    category_id?: string | "nope"
    cover_id?: string
    is_published?: boolean
    title: string
    content: string
    published_at?: string
    updated_at?: string

    // joined (GET-only)
    category_name?: string
    cover_path?: string
    cover_extension?: string
}