import { push } from 'svelte-spa-router'
import { wrap } from 'svelte-spa-router/wrap'
import { AuthStorage } from "@/tools/LocalStorage"
// views
import Index from '@/views/Index.svelte'
import Login from '@/views/Login.svelte'
import Logout from '@/views/Logout.svelte'
import Articles from '@/views/Articles.svelte'
import ArticleCreate from '@/views/ArticleCreate.svelte'
import Files from '@/views/Files.svelte'
import Settings from '@/views/Settings.svelte'


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
    //'*': Index, // other routes like 404
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