import ByteSize from 'byte-size'

export default class Size {

    /** convert bytes number to string like '10,4 MB' */
    public static convert(bytes: number): string {
        const data = ByteSize(bytes)
        return `${data.value} ${data.unit}`
    }

}