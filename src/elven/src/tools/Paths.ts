export class Env {
    
    // get uploads url
    public static getUploads(): string{
        return import.meta.env.VITE_UPLOADS_URL as string
    }

    // get api url
    public static getAPI(): string{
        return import.meta.env.VITE_API_URL as string
    }
    
}