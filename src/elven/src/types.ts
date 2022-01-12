import type { ElvenPlayer } from "@/plugins/elvenPlayer/types"
import type { ElvenNotify } from "@/plugins/elvenNotify/types"
import type { ElvenProgress } from "@/plugins/elvenProgress/types"

export type Meta = {
    per_page: number
    total_pages: number
    current_page: number
}

declare global {
    interface Window {
        $elvenPlayer: ElvenPlayer
        $elvenNotify: ElvenNotify
        $elvenProgress: ElvenProgress
    }
}