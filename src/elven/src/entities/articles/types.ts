import type { Meta } from "@/types"

export type Data = {
    meta: Meta
    data: Array<Article>
}

export enum Show {
    published = 'published',
    drafts = 'drafts'
}
export enum By {
    created = 'created',
    published = 'published', 
    updated = 'updated'
}
export enum Start {
    newest = 'newest',
    oldest = 'oldest'
}
export type Params = {
    page: number
    show: Show
    by: By
    start: Start
    preview: boolean
}

export type Article = {
    id?: string
    user_id?: string
    is_published?: boolean
    title: string
    content: string
    slug?: string
    published_at?: string
    updated_at?: string
}