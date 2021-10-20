//import type {theLogic as Player} from "@/common/plugins/ElvenPlayer/ElvenPlayer.js"
import type ElvenNotify from "@/common/plugins/ElvenNotify/ElvenNotify"
import type ElvenProgress from "@/common/plugins/ElvenProgress/ElvenProgress" 

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
        //$elvenPlayer: Player
        $elvenNotify: ElvenNotify
        $elvenProgress: ElvenProgress
    }
}