import type { Settings } from "./types";
import type { ElvenProgress as IElvenProgress } from "./types"

export default class ElvenProgress implements IElvenProgress {

    private line: HTMLDivElement
    private _percents: number
    // settings
    private settings: Settings = {
        height: '2px',
        basicLoading: {
            startTo: 45,
            startSpeed: 30,
            finishSpeed: 5
        }
    }

    constructor(line: HTMLDivElement, settings?: Settings) {
        if(settings) {
            this.settings = settings
        }
        this.line = line
        this.line.style.height = this.settings.height
        window.$elvenProgress = this
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
        const intervalID = setInterval(() => {
            if (this.percents < this.settings.basicLoading.startTo) {
                this.percents++
                return
            }
            clearInterval(intervalID)
        }, this.settings.basicLoading.startSpeed)
    }

    public finishBasic() {
        this.percents = this.settings.basicLoading.startTo
        const intervalID = setInterval(() => {
            if (this.percents < 100) {
                this.percents++
                return
            }
            clearInterval(intervalID)
            this.resetPercents()
        }, this.settings.basicLoading.finishSpeed)
    }


    public resetPercents() {
        this.percents = 0
    }
}