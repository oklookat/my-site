import type ElvenPlayerPlugin from "@/plugins/ElvenPlayer/ElvenPlayerPlugin"
import type { IElvenNotify } from "@/plugins/ElvenNotify/types"
import type { IElvenProgress } from "@/plugins/ElvenProgress/types"

export type TMeta = {
    per_page: number
    total_pages: number
    current_page: number
}

declare global {
    interface Window {
        $elvenPlayer: ElvenPlayerPlugin
        $elvenNotify: IElvenNotify
        $elvenProgress: IElvenProgress
    }
}