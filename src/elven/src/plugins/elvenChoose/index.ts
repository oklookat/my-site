import Confirm from "./confirm"

export interface IElvenChoose {
    confirm(title: string, warningText?: string): Promise<boolean>
}

export class ElvenChoose implements IElvenChoose {

    constructor() {
        this.init()
    }

    public init() {
        window.$choose = this
    }

    public async confirm(title: string, warningText?: string): Promise<boolean> {
        return new Confirm().start(title, warningText)
    }

}