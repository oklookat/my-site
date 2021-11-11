import type { exType } from '@/tools/Extensions'
import type { TMeta } from './global'

export interface TFilesData {
    meta: TMeta
    data: Array<TFile>
}

export interface TFile {
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
    extensionType?: exType
    sizeConverted?: string
    createdAtConverted?: string
}
