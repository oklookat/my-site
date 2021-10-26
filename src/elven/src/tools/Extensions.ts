export type exType = 'unknown' | 'image' | 'video' | 'audio'

export default class Extensions {

    private static readonly imageExtensions = ['jpeg', 'jpg', 'gif', 'png', 'bmp', 'svg', 'webp']
    private static readonly videoExtensions = ['mp4', 'mov', 'wmv', 'avi', 'flv', 'mkv', 'webm']
    private static readonly audioExtensions = ['mp3', 'm4a', 'flac', 'wav', 'wma', 'aac', 'alac', 'oog', 'aiff']

    public static getType(extension: string): exType {
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