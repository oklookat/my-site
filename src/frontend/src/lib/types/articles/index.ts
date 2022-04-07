/** request param */
export enum Show {
    published = 'published',
    drafts = 'drafts'
}

/** request param */
export enum By {
    created = 'created',
    published = 'published',
    updated = 'updated'
}

/** request param */
export enum Start {
    newest = 'newest',
    oldest = 'oldest'
}

/** request params */
export type Params = {
    page?: number
    show?: Show
    by?: By
    start?: Start
    /** full article or only title? */
    preview?: boolean
    //
    category_name?: string
    /** show articles without category? */
    without_category?: boolean
}

export function getDefaultParams(): Params {
    return {
        page: 1,
        show: Show.published,
        by: By.published,
        start: Start.newest,
        preview: true,
        category_name: null,
        without_category: false,
    }
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