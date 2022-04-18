const uploads = import.meta.env.VITE_UPLOADS_URL as string
const api = import.meta.env.VITE_API_URL as string

const uploadsURL = new URL(uploads)
const apiURL = new URL(api)

export class Env {

    /** get uploads URL (for uploading files etc) */
    public static getUploads(): URL {
        return uploadsURL
    }

    /**  get API URL */
    public static getAPI(): URL {
        return apiURL
    }

}

export class PathTools {

    /** api uploads + path. Like: 'https://uploads.example.com/yourpath' */
    public static getUploadsWith(path: string): URL {
        let final = path
        const uploads = Env.getUploads().toString()
        if (final.endsWith('/')) {
            final = `${uploads}${final}`
        } else {
            final = `${uploads}/${final}`
        }
        return new URL(final)
    }

    /** get path like: '/elven/yourpath' */
    public static getWithElvenPrefix(path: string): string {
        let pathd = path
        if (!pathd.startsWith('/')) {
            pathd = '/' + pathd
        }
        return `/elven${pathd}`
    }

}