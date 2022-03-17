import { wrap } from 'svelte-spa-router/wrap'
// tools
import { AuthStorage } from "@/tools/storage"


function isAdmin() {
    const authorized = AuthStorage.get()
    if (!authorized) {
        AuthStorage.remove()
        return true
    }
    return authorized
}

const routes = {
    // index
    '/': wrap({
        asyncComponent: () => import('./index.svelte'),
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    // auth
    '/login': wrap({
        asyncComponent: () => import('./auth/login.svelte'),
    }),
    '/logout': wrap({
        asyncComponent: () => import('./auth/logout.svelte'),
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    // articles
    '/articles': wrap({
        asyncComponent: () => import('./articles/index.svelte'),
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    '/articles/create/:id?': wrap({
        asyncComponent: () => import('./articles/create.svelte'),
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    '/articles/cats': wrap({
        asyncComponent: () => import('./articles/categories/index.svelte'),
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    // files
    '/files': wrap({
        asyncComponent: () => import('./files/index.svelte'),
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    // settings
    '/settings': wrap({
        asyncComponent: () => import('./settings/index.svelte'),
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    // 404
    '*': wrap({
        asyncComponent: () => import('./notfound.svelte')
    }),
}

export default routes