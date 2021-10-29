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
    onResponse: () => {
        console.log('response hooked')
        window.$elvenProgress.finish()
    },
    onError: (err) => {
        console.log(`${err.type} error hooked`)
        window.$elvenProgress.finish()
        AdapterError.handle(err)
    },
    onDownload: (e) => {
        console.log('download hooked')
    },
    onUploadProgress: (e) => {
        console.log('upload progress hooked')
        if (e.lengthComputable) {
            console.log(`Получено ${e.loaded} из ${e.total} байт`);
        } else {
            console.log(`Получено ${e.loaded} байт`); // если в ответе нет заголовка Content-Length
        }
    },
    onUploaded: () => {
        console.log('uploaded hooked')
    },
}

const config: TGlobalConfig = {
    withCredentials: true,
    baseURL: apiURL,
    hooks: hooks,
}

const Duck = new Ducksios(config)
export default Duck