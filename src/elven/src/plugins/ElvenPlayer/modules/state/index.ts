import type { Hooks, Hook, HookRemove } from "./types"
import { HookName } from "./types"

/** simple state management */
export default class State<T> {

    /** manipulated value */
    private value: T
    /** value hooks */
    private hooks: Hooks<T> = {
        onChange: {
            count: 0,
            items: {}
        }
    }

    constructor(value: T) {
        this.value = value
    }

    /** add user hook. Returns hook id */
    private addHook(name: HookName, hook: Hook<T>): number {
        const id = this.hooks[name].count++
        this.hooks[name].items[id] = hook
        return id
    }

    /** remove user hook */
    private removeHook(name: HookName, id: number) {
        const hookExists = id in this.hooks[name].items
        if (!hookExists) {
            return
        }
        delete this.hooks[name].items[id]
    }

    /** returns function that removes hook */
    private createHookRemover(name: HookName, id: number): HookRemove {
        const remove: HookRemove = () => {
            this.removeHook(name, id)
        }
        return remove
    }

    /** call callbacks in hook list */
    private notifyHooks(name: HookName, value: T) {
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
        this.notifyHooks(HookName.onChange, value)
    }

    /** add hook. Executes when value changed */
    public onChange(hook: Hook<T>): HookRemove {
        const id = this.addHook(HookName.onChange, hook)
        return this.createHookRemover(HookName.onChange, id)
    }

}