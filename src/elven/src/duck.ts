import Duckd, { type Config, type DuckHook } from "@oklookat/duck"

import { AdapterError } from "@/tools/errors"
import { Env } from "@/tools/paths"

const apiURL = Env.getAPI()

const Hooks: DuckHook.List = {
    onRequest(r) {
        window.$progress.startBasic()
    },
    onResponse() {
        window.$progress.finishBasic()
        //window.$progress.resetPercents()
    },
    onError(err) {
        window.$progress.finishBasic()
        window.$progress.reset()
        AdapterError.handle(err)
    },
    onUploadProgress(e) {
        // no Content-Length
        if (!e.data.lengthComputable) {
            return
        }
        const percents = (e.data.loaded / e.data.total) * 100
        window.$progress.percents = percents
    }
}


const config: Config = {
    timeout: 30000,
    withCredentials: true,
    baseURL: apiURL,
    hooks: Hooks,
}

const Duck = new Duckd(config)
export default Duck