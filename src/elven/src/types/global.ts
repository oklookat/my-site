import type { IElvenPlayer } from "@/plugins/ElvenPlayer/types"
import type { IElvenNotify } from "@/plugins/ElvenNotify/types"
import type { IElvenProgress } from "@/plugins/ElvenProgress/types"

export type TMeta = {
    per_page: number
    total_pages: number
    current_page: number
}

declare global {
    interface Window {
        $elvenPlayer: IElvenPlayer
        $elvenNotify: IElvenNotify
        $elvenProgress: IElvenProgress
    }
}