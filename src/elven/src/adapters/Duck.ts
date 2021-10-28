import { AdapterError } from "@/tools/ErrorHandler"
import { Env } from "@/tools/Paths"
import type { IHooks, TGlobalConfig } from "./ducksios/types"
import Ducksios from "./ducksios/general"

const apiURL = Env.getAPI()

const hooks: IHooks = {
    onRequest: () => {
        console.log(`request hooked`)
        window.$elvenProgress.start()
    },
    onRequestError: (err) => {
        console.log('request error hooked')
        window.$elvenProgress.finish()
        AdapterError.handle(err)
    },
    onResponse: () => {
        console.log('response hooked')
        window.$elvenProgress.finish()
    },
    onResponseError: (err) => {
        console.log('response error hooked')
        window.$elvenProgress.finish()
        AdapterError.handle(err)
    },
    onUploaded: () => {
        console.log('uploaded hooked')
    },
    onUploadError: (err) => {
        console.log('upload error hooked')
        window.$elvenProgress.finish()
        AdapterError.handle(err)
    },
    onUploadProgress: () => {
        console.log('upload progress hook')
    }
}

const config: TGlobalConfig = {
    withCredentials: true,
    baseURL: apiURL,
    hooks: hooks,
}

const Duck = new Ducksios(config)
export default Duck