import { By, type Params } from "$lib_elven/types/articles";

export default class ToolsArticles {

    /** get default request params */
    public static getDefaultParams(): Params {
        return {
            page: 1,
            published: true,
            newest: true,
            preview: true,
            by: By.published,
            title: null
        }
    }

    /** validate id */
    public static validateID(val: string): boolean {
        return !!val
    }

    /** validate title */
    public static validateTitle(val: string): boolean {
        return val.length <= 124
    }

    /** validate content */
    public static validateContent(val: string): boolean {
        return val && val.length <= 256000
    }

}