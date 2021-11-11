import type ElvenPlayerPlugin from "@/plugins/ElvenPlayer/ElvenPlayerPlugin"
import type ElvenNotifyPlugin from "@/plugins/ElvenNotify/ElvenNotifyPlugin"
import type { IElvenProgress } from "@/plugins/ElvenProgress/Types"

export type TMeta = {
    per_page: number
    total_pages: number
    current_page: number
}

declare global {
    interface Window {
        $elvenPlayer: ElvenPlayerPlugin
        $elvenNotify: ElvenNotifyPlugin
        $elvenProgress: IElvenProgress
    }
}