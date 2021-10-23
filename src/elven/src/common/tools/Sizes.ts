import ByteSize from 'byte-size'

export default class Sizes {
    static convert(bytes) {
        const data = ByteSize(bytes)
        return `${data.value} ${data.unit}`
    }
}