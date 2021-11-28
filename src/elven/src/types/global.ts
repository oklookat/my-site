import type { ElvenPlayer } from "@/plugins/ElvenPlayer/types"
import type { ElvenNotify } from "@/plugins/ElvenNotify/types"
import type { ElvenProgress } from "@/plugins/ElvenProgress/types"

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