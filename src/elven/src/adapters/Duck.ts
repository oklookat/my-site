import { AdapterError } from "@/tools/ErrorHandler"
import { Env } from "@/tools/Paths"
import type { IHooks, TGlobalConfig } from "@/plugins/ducksios/types"
import Ducksios from "@/plugins/ducksios/core"

const apiURL = Env.getAPI()


class Hooks implements IHooks {
    onRequest() {
        window.$elvenProgress.startBasic()
    }
    onResponse() {
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
    }
    onError(err) {
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
        AdapterError.handle(err)
    }
    onUploadProgress(e) {
        if (e.lengthComputable) {
            const percents = (e.loaded / e.total) * 100
            window.$elvenProgress.setPercents(percents)
            console.log(`uploaded: ${percents}%`)
        } else {
            // no Content-Length
            console.log(`uploaded ${e.loaded} bytes`);
        }
    }
    onUploaded() {
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
    }
}

const config: TGlobalConfig = {
    withCredentials: true,
    baseURL: apiURL,
    hooks: new Hooks(),
}

const Duck = new Ducksios(config)
export default Duck