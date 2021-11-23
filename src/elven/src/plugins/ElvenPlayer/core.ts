import type { IEvents, IState, TPlaylist, TSource, IElvenPlayer } from "./types";
import DOM from "./logic/dom";
import State from "./logic/state";
import Events from "./logic/events";
import Logger from "./logic/logger";
import type { Unsubscriber } from "svelte/store";


/** controls audio player */
export default class ElvenPlayer implements IElvenPlayer {

    public state: IState
    public dom: DOM
    private initialized: boolean = false
    private unsubs: Unsubscriber[] = []
    private events: IEvents
    private _playlist: TPlaylist = {
        position: 0,
        sources: []
    }

    constructor() {
        this.init()
    }

    public init() {
        this.state = new State()
        this.events = new Events(this.state)
        this.dom = new DOM(this.events)
        this.subscribe()
        this.initialized = true
        window.$elvenPlayer = this
    }

    public destroy() {
        this.unsubscribe()
        this.dom.destroy()
        this.dom = null
        this.initialized = false
        window.$elvenPlayer = null
    }

    private subscribe() {
        const u1 = this.state.store.current.ended.onChange(v => {
            if (v) {
                this.next()
            }
        })
        this.unsubs.push(u1)
    }

    private unsubscribe() {
        for (const unsub of this.unsubs) {
            unsub()
        }
        this.unsubs = []
    }

    public addToPlaylist(source: TSource) {
        this._playlist.sources.push(source)
    }

    public async play() {
        this.checkForInit()
        try {
            await this.dom.play()
        } catch (err) {
            Logger.error(err)
        }
    }

    public pause() {
        this.checkForInit()
        this.dom.pause()
    }

    public stop() {
        this.checkForInit()
        this.dom.stop()
    }

    private repeat() {
        this.stop()
        this.play()
    }

    public next() {
        const next = this.playlist.sources[this.playlist.position + 1]
        if (!next) {
            // stop if no source next
            this.stop()
            return
        }
        this.playlist.position++
        this.dom.source = next
        this.play()
    }

    public prev() {
        // if no source behind
        const prev = this.playlist.sources[this.playlist.position - 1]
        // if current time > 2% of total time - repeat
        const notInStart = this.dom.currentTimePercents > 2
        if (!prev || notInStart) {
            this.repeat()
            return
        }
        this.playlist.position--
        this.dom.source = prev
        this.play()
    }

    private checkForInit() {
        if (!this.initialized) {
            this.init()
        }
    }

    public get playlist(): TPlaylist {
        return this._playlist
    }

    public set playlist(playlist: TPlaylist) {
        const sources = playlist.sources
        if (sources.length === 0) {
            Logger.error('empty playlist')
            return
        }
        this._playlist = playlist
        let index = 0
        if (sources[playlist.position]) {
            index = playlist.position
        }
        this.dom.source = sources[index]
    }

}