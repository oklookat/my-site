import type { Meta } from "./global"

export type Data = {
    meta: Meta
    data: Array<Article>
}

export type Show = 'published' | 'drafts'
export type By = 'created' | 'published' | 'updated'
export type Start = 'newest' | 'oldest'
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
    content: Content
    slug?: string
    published_at?: string
    updated_at?: string
}

export type Content = {
    version?: string
    time?: number
    blocks: {
        id?: string
        type: any
        data: any
        //tunes?: {[name: string]: any}
    }[]
}