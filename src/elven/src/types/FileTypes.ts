import type { exType } from '@/tools/Extensions'
import type { IMeta } from './GlobalTypes'

export interface IFilesData {
    meta: IMeta
    data: Array<IFile>
}

export interface IFile {
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

export const IFileDefault: IFile = {
    id: '',
    user_id: '',
    hash: '',
    path: '',
    name: '',
    original_name: '',
    extension: '',
    size: 1,
    created_at: '',
    updated_at: '',
}
