export class AuthStorage {

    public static set(value: boolean) {
        const converted = value ? '1': '0'
        localStorage.setItem('elven-authenticated', converted)
    }

    public static get(): boolean {
        const state = localStorage.getItem('elven-authenticated')
        return state === '1'
    }

    public static remove() {
        localStorage.removeItem('elven-authenticated')
    }

}