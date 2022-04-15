import { invalidate } from "$app/navigation";
import Utils from "$lib_elven/tools";
import type { Data } from "$lib_elven/types";
import { By, Show, Start, type Article, type Params } from "$lib_elven/types/articles";

export default class ToolsArticles {

    public static getDefaultParams(): Params {
        return {
            page: 1,
            show: Show.published,
            by: By.published,
            start: Start.newest,
            preview: true,
            category_name: null,
            without_category: false,
        }
    }

    /** get items refresher */
    public static getRefresher(params: Params, items: Data<Article>, onPageChanged: (page: number) => Promise<void>): (isForce?: boolean) => Promise<void> {
        let force = false;
        let prevPage = params.page;
        let isFirstCall = true;
        const getData = async () => {
            if (isFirstCall && !force) {
                isFirstCall = false;
                return items.data;
            }
            if (prevPage < 2 && params.page < 2) {
                params.page = 1;
                await invalidate('');
                return items.data;
            }
            if (onPageChanged) {
                await onPageChanged(params.page);
            }
            return items.data;
        };
        const setPage = (val: number) => {
            prevPage = params.page;
            params.page = val;
        };
        return async (isForce = false) => {
            force = isForce;
            if (params.page < 2) {
                force = true
            }
            await Utils.refresh(params.page, setPage, getData);
        };
    }
}