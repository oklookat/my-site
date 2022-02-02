import "./style.scss"

enum className {
    container = "confirm",
    secondContainer = "confirm__second",
    title = "confirm__title",
    question = "confirm__question",
    NYContainer = "confirm__ny",
    no = "confirm__no",
    yes = "confirm__yes",
}

export default class Confirm {

    private container: HTMLDivElement
    private overflowSetter: () => void

    constructor() {
    }

    public async start(titleText: string, warningText?: string): Promise<boolean> {
        return new Promise(resolve => {
            // container
            const container = this.buildContainer(resolve)
            this.container = container
            // second container
            const secondContainer = this.buildSecondContainer()
            container.appendChild(secondContainer)
            // title
            const title = this.buildTitle(titleText)
            secondContainer.appendChild(title)
            // question
            const question = this.buildQuestion(warningText)
            secondContainer.appendChild(question)
            // no/yes container
            const NYC = this.buildNYContainer()
            secondContainer.appendChild(NYC)
            // no
            const no = this.buildNo(resolve)
            NYC.appendChild(no)
            // yes
            const yes = this.buildYes(resolve)
            NYC.appendChild(yes)
            //
            this.overflowSetter = this.getOverflowSetter()
            document.body.appendChild(this.container)
        })
    }


    private destroy() {
        this.overflowSetter()
        const isRemovable = this.container && this.container.parentNode === document.body
        if (!isRemovable) {
            return
        }
        document.body.removeChild(this.container)
    }

    /** first call sets overflow hidden, second call reset overflow to default */
    private getOverflowSetter(): () => void {
        const defOverflow = document.body.style.overflow
        document.body.style.overflow = 'hidden'
        return () => {
            document.body.style.overflow = defOverflow
        }
    }

    private buildContainer(resolve: (v: boolean) => void): HTMLDivElement {
        const container = document.createElement('div')
        container.className = className.container
        container.onclick = (e) => {
            e.stopPropagation()
            if (e.target !== container) {
                return
            }
            this.destroy()
            resolve(false)
        }
        return container
    }

    private buildSecondContainer(): HTMLDivElement {
        const container = document.createElement('div')
        container.className = className.secondContainer
        return container
    }

    private buildTitle(text: string): HTMLDivElement {
        const title = document.createElement('div')
        title.className = className.title
        title.innerText = text
        return title
    }

    private buildQuestion(warningText?: string): HTMLDivElement {
        const question = document.createElement('div')
        question.className = className.question
        if (warningText) {
            question.innerText = warningText
        } else {
            question.innerText = 'Are you sure?'
        }
        return question
    }

    private buildNYContainer(): HTMLDivElement {
        const NYC = document.createElement('div')
        NYC.className = className.NYContainer
        return NYC
    }

    private buildNo(resolve: (v: boolean) => void): HTMLDivElement {
        const no = document.createElement('div')
        no.className = className.no
        no.innerText = 'No'
        no.onclick = () => {
            this.destroy()
            resolve(false)
        }
        return no
    }

    private buildYes(resolve: (v: boolean) => void): HTMLDivElement {
        const yes = document.createElement('div')
        yes.className = className.yes
        yes.innerText = 'Yes'
        yes.onclick = () => {
            this.destroy()
            resolve(true)
        }
        return yes
    }

}