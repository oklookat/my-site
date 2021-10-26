import type ElvenPlayerPlugin from "@/plugins/ElvenPlayer/ElvenPlayerPlugin"
import type ElvenNotifyPlugin from "@/plugins/ElvenNotify/ElvenNotifyPlugin"
import type ElvenProgressPlugin from "@/plugins/ElvenProgress/ElvenProgressPlugin"

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