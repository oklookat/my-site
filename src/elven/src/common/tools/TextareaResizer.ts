let _global_this

export default class TextareaResizer {

    element: HTMLElement | null

    constructor(elementID: string) {
        this.element = document.getElementById(elementID)
        if(!this.element){
            throw Error(`TextareaResizer: element not found`)
        }
    }

    public start(){
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        _global_this = this
        onInput()
        // @ts-ignore
        this.element.addEventListener('input', onInput)
        window.addEventListener('resize', onInputByResize)
    }

    public destroy(){
        // @ts-ignore
        this.element.removeEventListener('input', onInput)
        window.removeEventListener('resize', onInputByResize)
    }
}


function onInputByResize(){
    setTimeout(() =>{
        onInput()
    }, 500)
}

function onInput(){
    _global_this.element.style.height = `0px`
    const height = _global_this.element.clientHeight
    const scrollHeight = _global_this.element.scrollHeight
    const isEqual = scrollHeight === height
    if(!isEqual){
        _global_this.element.style.height = `${scrollHeight}px`
    }
}