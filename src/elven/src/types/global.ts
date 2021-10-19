//import type {theLogic as Player} from "@/common/plugins/ElvenPlayer/ElvenPlayer.js"
//import type {theLogic as Notfiy} from "@/common/plugins/ElvenNotify/ElvenNotify.js"
//import type {theLogic as Progress} from "@/common/plugins/ElvenProgress/ElvenProgress.js"
 
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
        //$elvenNotify: Notfiy
        //$elvenProgress: Progress
    }
}