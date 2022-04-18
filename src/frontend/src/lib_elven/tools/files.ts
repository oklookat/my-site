import { invalidate } from "$app/navigation";
import Utils from "$lib_elven/tools";
import type { Data } from "$lib_elven/types";
import { By, Start, type File, type Params } from "$lib_elven/types/files";

export default class ToolsFiles {

    /** get default request params */
    public static getDefaultParams(): Params {
        return {
            page: 1,
            start: Start.newest,
            by: By.created,
            extensions: undefined,
            filename: undefined
        }
    }
}