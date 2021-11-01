import type { IMeta } from "./GlobalTypes"

export type TData = {
    meta: IMeta
    data: Array<TArticle>
}

export type TShow = 'published' | 'drafts'
export type TBy = 'created' | 'published' | 'updated'
export type TStart = 'newest' | 'oldest'
export type TParams = {
    page: number
    show: TShow
    by: TBy
    start: TStart
    preview: boolean
}

export type TArticle = {
    id?: string
    user_id?: string
    is_published?: boolean
    title: string
    content: TContent
    slug?: string
    published_at?: string
    updated_at?: string
}

export type TContent = {
    version?: string
    time?: number
    blocks: {
        id?: string
        type: any
        data: any
        tunes?: {[name: string]: any}
    }[]
}