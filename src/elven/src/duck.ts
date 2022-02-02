import Duckd, { type Cancelable } from "@oklookat/duck"
import type { Hooks as IHooks, GlobalConfig } from "@oklookat/duck"
import { AdapterError } from "@/tools/errors"
import { Env } from "@/tools/paths"

const apiURL = Env.getAPI()

export class CancelToken implements Cancelable {
    cancel(message?: string) {

    }
}

const Hooks: IHooks = {
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


const config: GlobalConfig = {
    timeout: 30000,
    withCredentials: true,
    baseURL: apiURL,
    hooks: Hooks,
}

const Duck = new Duckd(config)
export default Duck