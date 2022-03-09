/** article validator */
export default class Validate {

    /** validate id */
    public static id(val: string): boolean {
        return !!val
    }

    /** validate title */
    public static title(val: string): boolean {
        return val.length <= 124
    }

    /** validate content */
    public static content(val: string): boolean {
        return val && val.length <= 256000
    }

}