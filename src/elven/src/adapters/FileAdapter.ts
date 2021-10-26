import { Axios } from './Axios'
import type { IFilesData } from '@/types/FileTypes'

export default class FileAdapter {

    public static async getAll(page = 1, start = 'newest'): Promise<IFilesData> {
        const config = {
            params:
            {
                page: page, start: start
            }
        }
        try {
            const response = await Axios.get('files', config)
            return Promise.resolve(response.data as IFilesData)
        } catch (err) {
            return Promise.reject(err)
        }
    }

    public static async upload(files: FileList) {
        const config = { headers: { 'Content-Type': 'multipart/form-data' } }
        for (const file in files) {
            if (typeof files[file] !== "object") {
                continue
            }
            const formData = new FormData()
            formData.append("file", files[file])
            try {
                await Axios.post('files', formData, config)
            } catch (err) {
                return Promise.reject(err)
            }
        }
    }

    public static async delete(id: string) {
        try {
            await Axios.delete(`files/${id}`)
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err)
        }
    }
}