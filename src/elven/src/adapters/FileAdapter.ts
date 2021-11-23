import Duck from './Duck'
import type { TFilesData } from '@/types/file'

export default class FileAdapter {

    public static async getAll(page = 1, start = 'newest'): Promise<TFilesData> {
        const params = { page: page, start: start }
        try {
            const response = await Duck.GET({ url: 'files', params })
            return Promise.resolve(response.body as TFilesData)
        } catch (err) {
        }
    }

    public static async upload(files: FileList) {
        for (const file in files) {
            if (typeof files[file] !== 'object') {
                continue
            }
            const formData = new FormData()
            formData.append('file', files[file])
            try {
                await Duck.POST({ url: 'files', body: formData })
            } catch (err) {
            }
        }
    }

    public static async delete(id: string) {
        try {
            await Duck.DELETE({ url: `files/${id}` })
            return Promise.resolve()
        } catch (err) {
        }
    }
}