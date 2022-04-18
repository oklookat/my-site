import { By, type Params } from "$lib_elven/types/articles";

export default class ToolsArticles {

    /** get default request params */
    public static getDefaultParams(): Params {
        return {
            page: 1,
            published: true,
            newest: true,
            preview: true,
            without_category: false,
            by: By.published,
            category_name: null,
            title: null
        }
    }

}