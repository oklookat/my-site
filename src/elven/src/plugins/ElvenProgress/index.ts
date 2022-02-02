import type { Settings } from "./types";
import type { ElvenProgress as IElvenProgress } from "./types"

export default class ElvenProgress implements IElvenProgress {

    private basicActive = false
    private line: HTMLDivElement
    private _percents: number
    private settings: Settings = {
        height: '2px',
        basicLoading: {
            startTo: 45,
            startSpeed: 30,
            finishSpeed: 1
        }
    }

    constructor(line: HTMLDivElement, settings?: Settings) {
        if (settings) {
            this.settings = settings
        }
        this.line = line
        this.line.style.height = this.settings.height
        window.$progress = this
    }

    public get percents(): number {
        return this._percents
    }

    public set percents(percents: number) {
        if (percents > 100) {
            percents = 100
        } else if (percents < 0) {
            percents = 0
        }
        this._percents = percents
        this.line.style.width = `${percents}%`
    }

    public startBasic() {
        if (this.basicActive) {
            return
        }
        this.basicActive = true
        const intervalID = setInterval(() => {
            if (this.percents < this.settings.basicLoading.startTo) {
                this.percents++
                return
            }
            clearInterval(intervalID)
        }, this.settings.basicLoading.startSpeed)
    }

    public finishBasic() {
        if (!this.basicActive) {
            return
        }
        this.percents = this.settings.basicLoading.startTo
        const intervalID = setInterval(() => {
            if (this.percents < 100) {
                this.percents++
                return
            }
            this.reset()
            clearInterval(intervalID)
        }, this.settings.basicLoading.finishSpeed)
    }


    public reset() {
        this.basicActive = false
        this.percents = 0
    }
}