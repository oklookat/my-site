import Fetchd from '$lib/network'
import type { File as TFile, Params } from '$lib/types/files'
import type { Data } from '$lib/types'
import { StorageAuth } from '$lib/tools/storage'

export default class NetworkFile {

    private headers: Headers

    constructor(token: string) {
        const headers = new Headers()
        StorageAuth.addTokenToHeaders(headers, token)
        this.headers = headers
    }

    /** get files list */
    public async getAll(params: Params): Promise<Data<TFile>> {
        const paramsCopy: Params = { ...params }

        // convert params.extensionsSelector (if exists) to params.extensions
        if (paramsCopy.extensionsSelector) {
            let extensionsParsed: string[] = []

            // get selected extension to parse
            var extensions = paramsCopy.extensionsSelector.extensions
            var selected = paramsCopy.extensionsSelector.selected

            // if string - we need one type of file, like images
            if (typeof selected === "string") {
                // get one file types
                extensionsParsed = extensions[selected]
            } else if (selected instanceof Array) {
                // if array - we search many types of file, need concat that shit
                for (const readable of selected) {
                    const names = extensions[readable]
                    for (const extension of names) {
                        extensionsParsed.push(extension)
                    }
                }
            }
            // remove dups
            const extensionsUniq = [...new Set(extensionsParsed)];
            paramsCopy.extensions = extensionsUniq.join(",") as any
            params["extensions"] = paramsCopy.extensions
        }

        // send
        try {
            const response = await Fetchd.send({ method: "GET", url: 'files', params: paramsCopy, headers: this.headers })
            const jsond = await response.json()
            return jsond as Data<TFile>
        } catch (err) {
            return Promise.reject(err)
        }
    }

    /** upload one file */
    public async upload(file: File) {
        if (!(file instanceof File)) {
            return
        }
        const formData = new FormData()
        formData.append("file", file)
        try {
            await Fetchd.send({ method: "POST", url: 'files', body: formData, headers: this.headers })
            return
        } catch (err) {
            return Promise.reject(err)
        }
    }

    /** delete one file */
    public async delete(id: string) {
        try {
            await Fetchd.send({ method: "DELETE", url: `files/${id}`, headers: this.headers })
            return
        } catch (err) {
            return Promise.reject(err)
        }
    }
}