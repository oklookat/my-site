import Duckd from "@oklookat/duck"
import type { Hooks as IHooks, GlobalConfig } from "@oklookat/duck"
import { AdapterError } from "@/tools/errors"
import { Env } from "@/tools/paths"

const apiURL = Env.getAPI()

const Hooks: IHooks = {
    onRequest(r) {
        const method = r.config.method
        if (method === 'OPTIONS' || method === 'HEAD') {
            return
        }
        window.$elvenProgress.startBasic()
    },
    onResponse() {
        window.$elvenProgress.finishBasic()
        //window.$elvenProgress.resetPercents()
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
    }
}


const config: GlobalConfig = {
    withCredentials: true,
    baseURL: apiURL,
    hooks: Hooks,
}

const Duck = new Duckd(config)
export default Duck