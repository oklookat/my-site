import type { THooks, IHook, IHookRemove } from "./types"
import { THookName } from "./types"

/** basic state management */
export default class State<T> {

    /** manipulated value */
    private value: T
    /** value hooks */
    private hooks: THooks<T> = {
        onChange: {
            count: 0,
            items: {}
        }
    }

    constructor(value: T) {
        this.value = value
    }

    /** add user hook. Returns hook id */
    private addHook(name: THookName, hook: IHook<T>): number {
        const id = this.hooks[name].count++
        this.hooks[name].items[id] = hook
        return id
    }

    /** remove user hook */
    private removeHook(name: THookName, id: number) {
        const hookExists = id in this.hooks[name].items
        if (!hookExists) {
            return
        }
        delete this.hooks[name].items[id]
    }

    /** returns function that removes hook */
    private createHookRemover(name: THookName, id: number): IHookRemove {
        const remove: IHookRemove = () => {
            this.removeHook(name, id)
        }
        return remove
    }

    /** call callbacks in hook list */
    private notifyHooks(name: THookName, value: T) {
        const hooks = this.hooks[name].items
        for (const id in hooks) {
            hooks[id](value)
        }
    }

    /** get value */
    public get(): T {
        return this.value
    }

    /** set value */
    public set(value: T) {
        this.value = value
        this.notifyHooks(THookName.onChange, value)
    }

    /** add hook. Executes when value changed */
    public onChange(hook: IHook<T>): IHookRemove {
        const id = this.addHook(THookName.onChange, hook)
        return this.createHookRemover(THookName.onChange, id)
    }

}