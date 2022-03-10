import Duck from '@/network'
import type { File as TFile, Params } from '../entities/files/types'
import type { Data } from '@/types'

export default class FilesNetwork {

    /** get files list */
    public static async getAll(params: Params): Promise<Data<TFile>> {
        // format params.extensions to string like: "jpg,png,mp4,webp"
        if (params.extensions) {
            var selected = params.extensions.selected
            if (typeof selected === "string") {
                params.extensions = params.extensions.types[selected].join(",") as any
            } else if (selected instanceof Array) {
                const types = new Set()
                for (const selType of selected) {
                    types.add(params.extensions.types[selType])
                }
                // remove dups
                const typesUniq = [...new Set(types)];
                params.extensions = typesUniq.join(",") as any
            }
        }
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