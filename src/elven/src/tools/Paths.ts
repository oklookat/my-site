import { location } from "svelte-spa-router";

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

export class Route {

    /** returns location like: /articles */
    public static getLocation(): string {
        let _location = "";
        location.subscribe((val) => {
            _location = val;
        })();
        return _location;
    }

    /** get search params from window.location */
    public static getSearchParams(): URLSearchParams | undefined {
        const params = window.location.hash.split(/\?(.+)/)[1];
        if (!params) {
            return;
        }
        return new URLSearchParams(params);
    }


    /** set history query string by params */
    public static setHistoryParams(params: string | string[][] | 
        Record<string | number, string | number | boolean> | 
        URLSearchParams) {
        // @ts-ignore
        const _params = new URLSearchParams(params).toString();
        // set location
        const protocol = window.location.protocol + "//";
        const host = window.location.host;
        const pathname = window.location.pathname;
        const base = `${protocol}${host}${pathname}`;
        const _location = this.getLocation();
        const newurl = `${base}#${_location}?${_params}`;
        window.history.pushState({ params: _params }, "", newurl);
    }

}