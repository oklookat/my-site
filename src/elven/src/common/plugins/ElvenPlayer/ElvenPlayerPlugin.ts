import Player from './components/Player.svelte'
import PlayerLogic from "@/common/plugins/ElvenPlayer/components/PlayerLogic"


export default class ElvenPlayerPlugin {

    private player: Player
    private playerLogic: PlayerLogic

    constructor() {
        const el = document.querySelector('#elven__player')
        if (!el) {
            throw Error('elvenPlayer: element not found')
        }
        this.playerLogic = new PlayerLogic()
        this.player = new Player({
            target: el,
            props: {
                pl: this.playerLogic
            }
        })
        window.$elvenPlayer = this
    }

    public async addToPlaylist(url: string) {
        this.playerLogic.addToPlaylist(url)
    }

    public setPlaylist(playlist: string[]) {
        this.playerLogic.setPlaylist(playlist)
    }

    public async play(url: string) {
        await this.setPlaylist([url])
        await this.playerLogic.play()
        this.player.$set({ isActive: true })
    }

    public destroy() {
        this.player.$set({ isActive: false })
        this.playerLogic.destroy()
    }
}