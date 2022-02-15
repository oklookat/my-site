export type FileType = 'unknown' | 'image' | 'video' | 'audio'

export default class Extension {

    private static readonly imageExtensions = ['jpeg', 'jpg', 'gif', 'png', 'bmp', 'svg', 'webp']
    private static readonly videoExtensions = ['mp4', 'mov', 'wmv', 'avi', 'flv', 'mkv', 'webm']
    /** supported extensions to play audio in files */
    private static readonly audioExtensions = ['mp3', 'flac', 'wav', 'ogg']

    public static getType(extension: string): FileType {
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