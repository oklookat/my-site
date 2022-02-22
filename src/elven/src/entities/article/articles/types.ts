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
    preview?: boolean
    category_name?: string
    /** show articles without category */
    without_category?: boolean
}

/** article */
export type Article = {
    id?: string
    user_id?: string
    category_id?: string | "nope"
    category_name?: string
    is_published?: boolean
    title: string
    content: string
    published_at?: string
    updated_at?: string
}