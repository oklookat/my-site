import { push } from 'svelte-spa-router'
import { wrap } from 'svelte-spa-router/wrap'
import { AuthStorage } from "@/tools/localStorage"
// main
import Index from '@/entities/general/index.svelte'
// auth
import Login from '@/entities/auth/comps/login.svelte'
import Logout from '@/entities/auth/comps/logout.svelte'
// articles
import Articles from '@/entities/article/articles/comps/index.svelte'
import ArticleCreate from '@/entities/article/articles/comps/create.svelte'
import ArticlesCats from '@/entities/article/categories/comps/index.svelte'
// files
import Files from '@/entities/files/comps/index.svelte'
// settings
import Settings from '@/entities/settings/index.svelte'


function isAdmin() {
    const authorized = AuthStorage.get()
    if (!authorized) {
        push('/login')
        return true
    }
    return authorized
}

const routes = {
    '/': Index,
    // '*' = // TODO: add 404 route.
    '/login': wrap({
        component: Login,
    }),
    '/logout': wrap({
        component: Logout,
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    '/articles': wrap({
        component: Articles,
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    '/articles/create/:id?': wrap({
        component: ArticleCreate,
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    '/articles/cats': wrap({
        component: ArticlesCats,
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    '/files': wrap({
        component: Files,
        conditions: [
            () => { return isAdmin() }
        ]
    }),
    '/settings': wrap({
        component: Settings,
        conditions: [
            () => { return isAdmin() }
        ]
    }),
}

export default routes