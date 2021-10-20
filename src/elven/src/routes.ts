import {wrap} from 'svelte-spa-router/wrap'
import {push} from 'svelte-spa-router'
import { AuthStorage } from "@/common/tools/LStorage"
import Index from '@/views/Index.svelte'


function isAdmin() {
    const authorized = AuthStorage.get()
    if(!authorized){
        push('/login')
        return true
    }
    return authorized
}

const routes = {
    '/': Index,
    //'*': Index, // other routes like 404
    '/login': wrap({
        asyncComponent: () => import('@/views/auth/Login.svelte'),
    }),
    '/logout': wrap({
        asyncComponent: () => import('@/views/auth/Logout.svelte'),
        conditions: [
            () => {return isAdmin()}
        ]
    }),
    '/articles': wrap({
        asyncComponent: () => import('@/views/articles/Articles.svelte'),
        conditions: [
            () => {return isAdmin()}
        ]
    }),
    '/articles/create': wrap({
        asyncComponent: () => import('@/views/articles/ArticlesCreate.svelte'),
        conditions: [
            () => {return isAdmin()}
        ]
    }),
    '/files': wrap({
        asyncComponent: () => import('@/views/files/Files.svelte'),
        conditions: [
            () => {return isAdmin()}
        ]
    }),
    '/settings': wrap({
        asyncComponent: () => import('@/views/settings/Settings.svelte'),
        conditions: [
            () => {return isAdmin()}
        ]
    }),
}

export default routes