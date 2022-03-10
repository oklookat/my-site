/** generate file type selector by file type and extensions */
export function generateFileTypeSelector(readable: FileExtensionReadable | FileExtensionReadable[], 
    extensions: FileExtensions = FileExtensionsDefault): FileTypeSelector {
    const select: FileTypeSelector = {
        selected: readable,
        extensions: extensions
    }
    return select
}

export const FileExtensionsDefault: FileExtensions = {
    UNKNOWN: [],
    IMAGE: ['jpeg', 'jpg', 'gif', 'png', 'svg', 'bmp', 'webp'],
    VIDEO: ['mpg', 'mpeg', 'webm', 'mp4'],
    AUDIO: ['mp3', 'flac', 'wav', 'ogg']
}

export interface FileTypeSelector {
    selected: FileExtensionReadable | FileExtensionReadable[]
    get extensions(): FileExtensions
}

export type FileExtensionReadable = 'UNKNOWN' | 'IMAGE' | 'VIDEO' | 'AUDIO'

export interface FileExtensions {
    UNKNOWN: string[]
    /** browser must support this image extensions */
    IMAGE: string[]
    /** browser must support this video extensions */
    VIDEO: string[]
    /** browser must support this audio extensions */
    AUDIO: string[]
}

export default class Extension {

    /** get readable file type.
     * @param extension extension without dot
     */
    public static getType(extension?: string): FileTypeSelector {
        const select = generateFileTypeSelector("UNKNOWN")
        if (!extension) {
            return select
        }
        extension = extension.toLowerCase()
        const image = select.extensions.IMAGE.includes(extension)
        if (image) {
            select.selected = "IMAGE"
            return select
        }
        const video = select.extensions.VIDEO.includes(extension)
        if (video) {
            select.selected = "VIDEO"
            return select
        }
        const audio = select.extensions.AUDIO.includes(extension)
        if (audio) {
            select.selected = "AUDIO"
            return select
        }
        return select
    }

}