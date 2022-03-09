/** generate file type selector */
export function generateFileTypeSelector(theType: FileType | FileType[], types: FileTypes = FileTypeDefault): FileTypeSelector {
    const select: FileTypeSelector = {
        selected: theType,
        types: types
    }
    return select
}

export const FileTypeDefault: FileTypes = {
    UNKNOWN: [],
    IMAGE: ['jpeg', 'jpg', 'gif', 'png', 'svg', 'bmp', 'webp'],
    VIDEO: ['mpg', 'mpeg', 'webm', 'mp4'],
    AUDIO: ['mp3', 'flac', 'wav', 'ogg']
}

export interface FileTypeSelector {
    selected: FileType | FileType[]
    get types(): FileTypes
}

export type FileType = 'UNKNOWN' | 'IMAGE' | 'VIDEO' | 'AUDIO'
export interface FileTypes {
    UNKNOWN: string[]
    /** browser must support this image extensions */
    IMAGE: string[]
    /** browser must support this video extensions */
    VIDEO: string[]
    /** browser must support this audio extensions */
    AUDIO: string[]
}

export default class Extension {

    public static getType(extension?: string): FileTypeSelector {
        const select = generateFileTypeSelector("UNKNOWN")
        if (!extension) {
            return select
        }
        extension = extension.toLowerCase()
        const image = select.types.IMAGE.includes(extension)
        if (image) {
            select.selected = "IMAGE"
            return select
        }
        const video = select.types.VIDEO.includes(extension)
        if (video) {
            select.selected = "VIDEO"
            return select
        }
        const audio = select.types.AUDIO.includes(extension)
        if (audio) {
            select.selected = "AUDIO"
            return select
        }
        return select
    }

}