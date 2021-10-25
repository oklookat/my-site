import { wrap } from 'svelte-spa-router/wrap'
import { push } from 'svelte-spa-router'
import { AuthStorage } from "@/common/tools/LStorage"

import Index from '@/views/Index.svelte'
import Login from '@/views/auth/Login.svelte'
import Logout from '@/views/auth/Logout.svelte'
import Articles from '@/views/articles/Articles.svelte'
import ArticleCreate from '@/views/articles/ArticleCreate.svelte'
import Files from '@/views/files/Files.svelte'
import Settings from '@/views/settings/Settings.svelte'

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