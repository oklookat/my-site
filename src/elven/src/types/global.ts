import type ElvenPlayerPlugin from "@/common/plugins/ElvenPlayer/ElvenPlayerPlugin"
import type ElvenNotifyPlugin from "@/common/plugins/ElvenNotify/ElvenNotifyPlugin"
import type ElvenProgressPlugin from "@/common/plugins/ElvenProgress/ElvenProgressPlugin"

export interface IMeta {
    per_page: number
    total_pages: number
    current_page: number
}

export const iMetaDefault: IMeta = {
    per_page: 0,
    total_pages: 0,
    current_page: 0
}

declare global {
    interface Window {
        $elvenPlayer: ElvenPlayerPlugin
        $elvenNotify: ElvenNotifyPlugin
        $elvenProgress: ElvenProgressPlugin
    }
}