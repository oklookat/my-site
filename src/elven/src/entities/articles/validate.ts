import type { Params } from "./types";
import { Show, By, Start } from "./types";

export default class Validate {

    /** validate and correct URLSearchParams by Params */
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
        const hasShow = urlParams.has("show");
        if (hasShow) {
            const show = urlParams.get("show");
            if (show in Show) {
                params.show = Show[show];
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
        const hasPreview = urlParams.has("preview");
        if (hasPreview) {
            const preview = urlParams.get("preview").toLowerCase();
            let parsed: boolean | undefined = undefined;
            if (preview === "true") {
                parsed = true;
            } else if (preview === "false") {
                parsed = false;
            }
            if (parsed !== undefined) {
                params.preview = parsed;
            }
        }
    }

}