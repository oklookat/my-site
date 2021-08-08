export default class Extensions {

    static getReadable(extension: string): string {
        extension = extension.toUpperCase()
        const isImage = extension === 'JPEG' || extension === 'JPG' ||
            extension === 'GIF' || extension === 'PNG' || extension === 'BMP' || extension === 'SVG' || extension === 'WEBP'
        if (isImage) {
            return 'IMAGE'
        }
        const isVideo = extension === 'MP4' || extension === 'MOV' || extension === 'WMV'
            || extension === 'AVI' || extension === 'FLV' || extension === 'MKV' || extension === 'WEBM'
        if (isVideo) {
            return 'VIDEO'
        }
        const isAudio = extension === 'MP3' || extension === 'M4A' || extension === 'FLAC'
            || extension === 'WAV' || extension === 'WMA' || extension === 'AAC' || extension === 'ALAC' || extension === 'OOG'
            || extension === 'AIFF'
        if (isAudio) {
            return 'AUDIO'
        }
        return 'UNKNOWN'
    }

}