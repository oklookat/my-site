import theLogic from "../common/plugins/ElvenPlayer/ElvenPlayer.js"

export interface IMeta {
    per_page: number
    next: string
}

export const iMetaDefault: IMeta = {
    per_page: 0,
    next: '',
}

declare global {
    interface Window {
        $elvenPlayer: theLogic;
    }
}