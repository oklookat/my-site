import { AdapterError } from "@/tools/ErrorHandler"
import { Env } from "@/tools/Paths"
import type { Hooks as IHooks, GlobalConfig } from "@/plugins/ducksios/types"
import Ducksios from "@/plugins/ducksios"

const apiURL = Env.getAPI()

const Hooks: IHooks = {
    onRequest() {
        window.$elvenProgress.startBasic()
    },
    onResponse() {
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
    },
    onError(err) {
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
        AdapterError.handle(err)
    },
    onUploadProgress(e) {
        if (!e.data.lengthComputable) {
            // no Content-Length
            console.log(`uploaded ${e.data.loaded} bytes`);
            return
        }
        const percents = (e.data.loaded / e.data.total) * 100
        window.$elvenProgress.percents = percents
        console.log(`uploaded: ${percents}%`)
    },
    onUploaded() {
        window.$elvenProgress.finishBasic()
        window.$elvenProgress.resetPercents()
    }
}


const config: GlobalConfig = {
    withCredentials: true,
    baseURL: apiURL,
    hooks: Hooks,
}

const Duck = new Ducksios(config)
export default Duck