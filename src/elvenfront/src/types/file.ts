import { IMeta } from '@/types/global'

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
