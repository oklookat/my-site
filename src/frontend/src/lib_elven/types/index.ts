import type { ElvenPlayer } from "$lib_elven/plugins/elvenPlayer/types"
import type { ElvenNotify } from "$lib_elven/plugins/elvenNotify/types"
import type { ElvenProgress } from "$lib_elven/plugins/elvenProgress/types"
import type { IElvenChoose } from "../plugins/elvenChoose"

export type Page = number
export type Counter = number | string

/** response with multiple entities */
export type Items<T> = {
    meta: Meta
    data: Record<Counter, T>
}

/** information about requested data */
export type Meta = {
    per_page: Page
    total_pages: Page
    current_page: Page
}

declare global {
    interface Window {
        // plugins
        $player: ElvenPlayer
        $notify: ElvenNotify
        $progress: ElvenProgress
        $choose: IElvenChoose
    }
}