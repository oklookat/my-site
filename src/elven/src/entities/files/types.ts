import type { FileType } from '@/tools/extension'

/** request param */
export enum Start {
    newest = 'newest',
    oldest = 'oldest'
}

/** request param */
export enum By {
    created = 'created'
}

/** request params */
export type Params = {
    page: number
    start: Start
    by: By
}

/** file */
export type File = {
    id: string
    user_id: string
    hash: string
    path: string
    name: string
    original_name: string
    extension: string
    size: number
    created_at: string
    updated_at: string
    // not in model
    pathConverted?: URL
    sizeConverted?: string
    createdAtConverted?: string
    extensionType?: FileType
}
