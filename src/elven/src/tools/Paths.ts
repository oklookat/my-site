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