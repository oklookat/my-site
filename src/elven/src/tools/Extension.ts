export type FileType = 'unknown' | 'image' | 'video' | 'audio'

export default class Extension {

    /** browser support this image extensions */
    private static readonly imageExtensions = ['jpeg', 'jpg', 'gif', 'png', 'svg', 'bmp', 'webp']
    /** browser support this video extensions */
    private static readonly videoExtensions = ['mpg', 'mpeg', 'webm', 'mp4']
    /** browser support this audio extensions */
    private static readonly audioExtensions = ['mp3', 'flac', 'wav', 'ogg']

    public static getType(extension?: string): FileType {
        if (!extension) {
            return "unknown"
        }
        extension = extension.toLowerCase()
        const image = this.imageExtensions.includes(extension)
        if (image) {
            return 'image'
        }
        const video = this.videoExtensions.includes(extension)
        if (video) {
            return 'video'
        }
        const audio = this.audioExtensions.includes(extension)
        if (audio) {
            return 'audio'
        }
        return 'unknown'
    }

}