export class AuthStorage {

    public static set(value: boolean) {
        let converted = "0"
        if (value) {
            converted = "1"
        }
        localStorage.setItem("elven-authenticated", converted)
    }

    public static get(): boolean {
        const state = localStorage.getItem("elven-authenticated")
        return state === "1"
    }

    public static remove() {
        localStorage.removeItem("elven-authenticated")
    }

}