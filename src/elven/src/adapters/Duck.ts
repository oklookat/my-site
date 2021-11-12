import { AdapterError } from "@/tools/ErrorHandler"
import { Env } from "@/tools/Paths"
import type { IHooks, TGlobalConfig } from "@/plugins/ducksios/types"
import Ducksios from "@/plugins/ducksios/core"

const apiURL = Env.getAPI()


const hooks: IHooks = {
    onRequest: () => {
        console.log('on request hooked')
        window.$elvenProgress.startBasic()
    },
    onResponse: () => {
        console.log('on response hooked')
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
    },
    onError: (err) => {
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
        AdapterError.handle(err)
    },
    onUploadProgress: (e) => {
        console.log('upload progress hooked')
        if (e.lengthComputable) {
            const percents = (e.loaded / e.total) * 100
            window.$elvenProgress.setPercents(percents)
            console.log(`uploaded: ${percents}%`)
        } else {
            // no Content-Length
            console.log(`uploaded ${e.loaded} bytes`);
        }
    },
    onUploaded: () => {
        console.log('uploaded hooked')
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
    },
}

const config: TGlobalConfig = {
    withCredentials: true,
    baseURL: apiURL,
    hooks: hooks,
}

const Duck = new Ducksios(config)
export default Duck