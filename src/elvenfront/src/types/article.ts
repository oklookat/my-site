import { IMeta } from "./global"

export interface IArticlesData {
    meta: IMeta
    data: Array<IArticle>
}

export interface IArticle {
    id: string
    user_id: string
    is_published: boolean
    title: string
    content: articleContent
    slug: string
    published_at: string
    updated_at: string
}

export type articleContent = {
	time: number
	blocks: {
		id: string
		type: string
		data: any
	}
	version: string
}