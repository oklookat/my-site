import type { TFileType } from '@/tools/Extensions'
import type { TMeta } from './global'

export type TFilesData = {
    meta: TMeta
    data: Array<TFile>
}

export type TFile = {
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
    //
    extensionType?: TFileType
    sizeConverted?: string
    createdAtConverted?: string
}
