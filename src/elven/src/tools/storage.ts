import { push } from 'svelte-spa-router';
import { writable } from 'svelte/store';

/** app state */
export class GlobalState {

    /** is user now in 404 page? */
    public static isNotFoundPage = writable(false);

}

/** auth storage */
export class AuthStorage {

    /** set user login state */
    public static set(value: boolean) {
        const converted = value ? '1' : '0'
        localStorage.setItem('elven-authenticated', converted)
    }

    /** is user logged in? (not trust this value, its localstorage) */
    public static get(): boolean {
        const state = localStorage.getItem('elven-authenticated')
        return state === '1'
    }

    /** remove auth state from storage & push to login route */
    public static remove() {
        localStorage.removeItem('elven-authenticated')
        push("/login")
    }

}