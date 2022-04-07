import type { FileTypeSelector } from '../../tools/extension'

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
    page?: number
    start?: Start
    by?: By
    /** find files with this extensions. Format like: 'jpg,gif,png' */
    extensions?: string
    filename?: string
    // not in model
    extensionsSelector?: FileTypeSelector
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
    extensionsSelector: FileTypeSelector
}
