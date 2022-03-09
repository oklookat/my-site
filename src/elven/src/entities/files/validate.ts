import type { Params } from "./types";
import { By, Start } from "./types";

export default class Validate {

    /** validate and correct Params by URLSearchParams */
    public static params(params: Params, urlParams?: URLSearchParams) {
        if (!urlParams) {
            return;
        }
        const hasPage = urlParams.has("page");
        if (hasPage) {
            const page = urlParams.get("page");
            const parsed = parseInt(page, 10);
            if (!isNaN(parsed)) {
                params.page = parsed;
            }
        }
        const hasBy = urlParams.has("by");
        if (hasBy) {
            const by = urlParams.get("by");
            if (by in By) {
                params.by = By[by];
            }
        }
        const hasStart = urlParams.has("start");
        if (hasStart) {
            const start = urlParams.get("start");
            if (start in Start) {
                params.start = Start[start];
            }
        }
    }

}