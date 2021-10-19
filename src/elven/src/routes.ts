import { AuthStorage } from "@/common/tools/LStorage"
import Authorized from "@/layouts/Authorized.svelte"
import NotAuthorized from "@/layouts/NotAuthorized.svelte"
import Index from '@/views/Index.svelte'
import Login from '@/views/auth/Login.svelte'
import Logout from '@/views/auth/Logout.svelte'
import Articles from '@/views/articles/Articles.svelte'
import ArticlesCreate from '@/views/articles/ArticleCreate.svelte'
import Files from '@/views/files/Files.svelte'
import Settings from '@/views/settings/Settings.svelte'

function isAdmin() {
    const authorized = AuthStorage.get()
    return authorized
}

const routes = [
    {
        name: '/',
        component: Index,
        layout: Authorized,
        onlyIf: { guard: isAdmin, redirect: '/login' }
    },
    { name: '/login', component: Login, layout: NotAuthorized },
    { name: '/logout', component: Logout, layout: Authorized, onlyIf: { guard: isAdmin, redirect: '/login' } },
    { name: '/articles', component: Articles, layout: Authorized, onlyIf: { guard: isAdmin, redirect: '/login' } },
    { name: '/articles/create', component: ArticlesCreate, layout: Authorized, onlyIf: { guard: isAdmin, redirect: '/login' } },
    { name: '/files', component: Files, layout: Authorized, onlyIf: { guard: isAdmin, redirect: '/login' } },
    { name: '/settings', component: Settings, layout: Authorized, onlyIf: { guard: isAdmin, redirect: '/login' } },
]

export { routes }