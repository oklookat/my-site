import type { ElvenPlayer } from "@/plugins/elvenPlayer/types"
import type { ElvenNotify } from "@/plugins/elvenNotify/types"
import type { ElvenProgress } from "@/plugins/elvenProgress/types"
import type { IElvenChoose } from "./plugins/elvenChoose"

/** response with multiple entities */
export type Data<T> = {
    meta: Meta
    data: Record<number, T>
}

/** information about requested data */
export type Meta = {
    per_page: number
    total_pages: number
    current_page: number
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