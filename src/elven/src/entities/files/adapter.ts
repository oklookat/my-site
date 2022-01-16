import Duck from '@/duck'
import type { File as TFile, Params } from './types'
import type { Data } from '@/types'

export default class FileAdapter {

    /** get files list */
    public static async getAll(params: Params): Promise<Data<TFile>> {
        try {
            const response = await Duck.GET({ url: 'files', params })
            return Promise.resolve(response.body as Data<TFile>)
        } catch (err) {
        }
    }

    /** upload one file */
    public static async upload(file: File) {
        if (!(file instanceof File)) {
            return
        }
        const formData = new FormData()
        formData.append("file", file)
        try {
            await Duck.POST({ url: 'files', body: formData })
            return Promise.resolve()
        } catch (err) {
        }
    }

    /** delete one file */
    public static async delete(id: string) {
        try {
            await Duck.DELETE({ url: `files/${id}` })
            return Promise.resolve()
        } catch (err) {
        }
    }
}