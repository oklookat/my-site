const uploads = import.meta.env.VITE_UPLOADS_URL as string
const api = import.meta.env.VITE_API_URL as string

export class Env {

    /** get uploads URL (for uploading files etc) */
    public static getUploads(): string {
        return uploads
    }

    /**  get API URL */
    public static getAPI(): string {
        return api
    }

}

export class PathTools {

    /** get uploads URL + path */
    public static getUploadsWith(path: string): URL {
        const final = `${Env.getUploads()}/${path}`;
        return new URL(final)
    }

    public static getWithElvenPrefix(path: string): string {
        return `/elven${path}`
    }

}

export class Route {

    /** init onpopstate and set  */
    public static initPopState(callback: (searchParams: URLSearchParams) => void) {
        const startLocation = Route.getLocation()
        // back/forward buttons in browser triggered / location changed
        window.onpopstate = () => {
            const newLocation = Route.getLocation()
            if (newLocation !== startLocation) {
                // destroy if location changed
                window.onpopstate = undefined;
                return;
            }
            const searchParams = Route.getSearchParams();
            callback(searchParams)
        };
    }

    /** 
    * [0] = full path. Ex: #/articles?params
    * 
    * [1] = ex: /articles
    * 
    * [2] = params with ? in start. Ex: ?page=1&show=published&by=updated&start=newest&preview=true
    *
     */
    public static getLocation(): string {
        const regexp = /^#(\/.*?)(\?.*)?$/gm
        const result = regexp.exec(window.location.hash)
        const path = result[1]
        if (!path) {
            return "/"
        }
        return path
    }

    /** get search params from window.location */
    public static getSearchParams(): URLSearchParams | undefined {
        const params = window.location.hash.split(/\?(.+)/)[1]
        if (!params) {
            return
        }
        return new URLSearchParams(params)
    }


    /** set history query string by params */
    public static setHistoryParams(params: string | string[][] |
        Record<string | number, string | number | boolean> |
        URLSearchParams | null) {
        let _params = ""
        if(params) {
            // @ts-ignore
            _params = "?" + new URLSearchParams(params).toString()
        }
        // set location
        const protocol = window.location.protocol + "//"
        const host = window.location.host;
        const pathname = window.location.pathname
        const base = `${protocol}${host}${pathname}`
        const _location = this.getLocation()
        const newurl = `${base}#${_location}${_params}`
        window.history.pushState({ params: _params }, "", newurl)
    }

}