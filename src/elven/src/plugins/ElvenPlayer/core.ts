import DOM from "./logic/dom";
import type { IEvents, IState, TPlaylist, TSource } from "./types";
import State from "./logic/state";
import Events from "./logic/events";


export default class Core {

    public state: IState
    private events: IEvents
    public dom: DOM
    private _playlist: TPlaylist = {
        position: 0,
        sources: []
    }

    constructor() {
        this.state = new State()
        this.events = new Events(this.state, this.next)
        this.dom = new DOM('', this.events)
    }

    public destroy() {
        this.dom.destroy()
    }

    public get playlist(): TPlaylist {
        return this._playlist
    }

    public set playlist(playlist: TPlaylist) {
        const sources = playlist.sources
        if (sources.length < 1) {
            throw Error('ElvenPlayer: empty playlist.')
        }
        this._playlist = playlist
        let index = 0
        if(sources[playlist.position]) {
            index = playlist.position
        }
        this.dom.source = sources[index]
    }

    public addToPlaylist(source: TSource) {
        this._playlist.sources.push(source)
    }

    public async play() {
        await this.dom.play()
    }

    public pause() {
        this.dom.pause()
    }

    public stop() {
        this.dom.stop()
    }

    public next() {
        const next = this.playlist[this.playlist.position + 1]
        if (!next) {
            this.stop()
            return
        }
        this.playlist.position++
        this.dom.source = next
        this.play()
    }

    public prev() {
        const prev = this.playlist[this.playlist.position - 1]
        if (!prev) {
            // repeat
            this.stop()
            this.play()
            return
        }
        this.playlist.position--
        this.dom.source = prev
        this.play()
    }

}