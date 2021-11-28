import type { FileType } from '@/tools/Extensions'
import type { Meta } from './global'

export type Data = {
    meta: Meta
    data: Array<File>
}

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
