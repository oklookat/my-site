/** wtf is going on:
 * 1. We need to search files by extensions / or after getting file from server we need readable extension
 * 2. API recieves param for searching extensions in format like: 'jpg,png,mp4' etc
 * 3. To easily manipulate this we need fileTypeSelector
 * 4. FTS contains extensions for file types, and selected file type(s)
 * 5. Selected = what extensions we selected to search / extensions = FileTypeSelector.extensions[selected]
 */


/** generate file type selector by file type and extensions */
export function generateFileTypeSelector(readable: FileExtensionReadable | FileExtensionReadable[],
    extensions: FileExtensions = FileExtensionsDefault): FileTypeSelector {
    if (readable instanceof Array) {
        // remove dups
        readable = [...new Set(readable)]
    }
    const select: FileTypeSelector = {
        selected: readable,
        extensions: extensions,
        selectedToString: () => { return selectedToString(readable, extensions) }
    }
    return select
}

function selectedToString(readable: Readable, extensions: FileExtensions = FileExtensionsDefault): string {
    if (!(readable instanceof Array)) {
        return extensions[readable].join(',')
    }
    const extsStrings = []
    for (const read of readable) {
        const exts = extensions[read]
        extsStrings.push(exts.join(','))
    }
    return extsStrings.join(',')
}

export type Readable = FileExtensionReadable | FileExtensionReadable[]

export type FileExtensionReadable = 'UNKNOWN' | 'IMAGE' | 'VIDEO' | 'AUDIO'

export const FileExtensionsDefault: FileExtensions = {
    UNKNOWN: [],
    IMAGE: ['jpeg', 'jpg', 'gif', 'png', 'svg', 'bmp', 'webp'],
    VIDEO: ['mpg', 'mpeg', 'webm', 'mp4'],
    AUDIO: ['mp3', 'flac', 'wav', 'ogg']
}

export interface FileTypeSelector {
    /** file type */
    selected: FileExtensionReadable | FileExtensionReadable[]

    /** selected file extension(s) to string */
    selectedToString: () => string

    /** file extensions by types */
    get extensions(): FileExtensions
}


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

    /** get file type selector by extension.
     * @param extension extension without dot
     */
    public static getSelector(extension?: string): FileTypeSelector {
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