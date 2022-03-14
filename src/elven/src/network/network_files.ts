import Duck from '@/network'
import type { File as TFile, Params } from '@/types/files'
import type { Data } from '@/types'

export default class NetworkFile {

    /** get files list */
    public static async getAll(params: Params): Promise<Data<TFile>> {
        const paramsCopy: Params = {...params}
        // convert params.extensionsSelector to params.extensions
        if (paramsCopy.extensionsSelector) {
            let extensionsParsed: string[] = []

            // get selected extension to parse
            var extensions = paramsCopy.extensionsSelector.extensions
            var selected = paramsCopy.extensionsSelector.selected

            // if string - we need one type of file, like images
            if(typeof selected === "string") {
                // get one file types
                extensionsParsed = extensions[selected]
            } else if(selected instanceof Array) {
                // if array - we search many types of file, need concat that shit
                for(const readable of selected) {
                    const names = extensions[readable]
                    for(const extension of names) {
                        extensionsParsed.push(extension)
                    }
                }
            }
            // remove dups
            const extensionsUniq = [...new Set(extensionsParsed)];
            paramsCopy.extensions = extensionsUniq.join(",") as any
            params["extensions"] = paramsCopy.extensions
        }
        try {
            const response = await Duck.GET({ url: 'files', params: paramsCopy })
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