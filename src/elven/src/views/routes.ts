import { push } from 'svelte-spa-router'
import { wrap } from 'svelte-spa-router/wrap'
// tools
import { AuthStorage } from "@/tools/storage"
// // main
// import Index from './index.svelte'
// // auth
// import Login from './auth/login.svelte'
// import Logout from './auth/logout.svelte'
// // articles
// import Articles from './articles/index.svelte'
// import ArticleCreate from './articles/create.svelte'
// import ArticlesCats from './articles/categories/index.svelte'
// // files
// import Files from './files/index.svelte'
// // settings
// import Settings from './settings/index.svelte'


function isAdmin() {
    const authorized = AuthStorage.get()
    if (!authorized) {
        push('/login')
        return true
    }
    return authorized
}

const routes = {
    // index
    '/': wrap({
        asyncComponent: () => import('./index.svelte')
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