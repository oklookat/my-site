//
import Fetchd from '$lib_elven/network'
import { addTokenToHeaders } from '$lib_elven/tools'
import type { Params } from '$lib_elven/types/files'

/** Use with SSR by passing token / or in components by passing empty token.
 * 
 * Static methods = not for SSR, use only on components side
 */
export default class NetworkFile {

    private headers: Headers

    constructor(token: string) {
        const headers = new Headers()
        addTokenToHeaders(token, headers)
        this.headers = headers
    }

    /** get files list */
    public async getAll(params: Params): Promise<Response> {
        // send
        try {
            const response = await Fetchd.send({ method: "GET", url: 'files', params: params, headers: this.headers })
            return response
        } catch (err) {
            throw err
        }
    }

    /** upload one file */
    public static async upload(file: File): Promise<Response> {
        if (!(file instanceof File)) {
            return
        }
        const form = new FormData()
        //
        form.append("file", file)
        try {
            const resp = await Fetchd.send({ method: "POST", url: 'files', body: form })
            return resp
        } catch (err) {
            return Promise.reject(err)
        }
    }

    /** delete one file */
    public static async delete(id: string) {
        try {
            await Fetchd.send({ method: "DELETE", url: `files/${id}` })
            return
        } catch (err) {
            return Promise.reject(err)
        }
    }
}