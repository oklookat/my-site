import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router'
import Store from "@/store/index"
import Login from '@/views/Auth/Login.vue'
import Logout from '@/views/Auth/Logout.vue'
import Index from '@/views/Main/Index.vue'
import Articles from '@/views/Main/Articles/Articles.vue'
import ArticleCreate from '@/views/Main/Articles/ArticleCreate.vue'
import Files from '@/views/Main/Files/Files.vue'
import Settings from '@/views/Main/Settings/Settings.vue'
import AuthAdapter from "@/common/adapters/Main/AuthAdapter";
import Fetcher from "@/common/adapters/Fetcher";

const routes: Array<RouteRecordRaw> = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: {auth: 'no'}
    },
    {
        path: '/logout',
        name: 'Logout',
        component: Logout,
        meta: {auth: 'yes'}
    },
    {
        path: '/',
        name: 'Index',
        component: Index,
        meta: {auth: 'yes'}
    },
    {
        path: '/articles',
        name: 'Articles',
        component: Articles,
        meta: {auth: 'yes'}
    },
    {
        path: '/articles/create/:id?',
        name: 'ArticleCreate',
        component: ArticleCreate,
        meta: {auth: 'yes'}
    },
    {
        path: '/files',
        name: 'Files',
        component: Files,
        meta: {auth: 'yes'}
    },
    {
        path: '/settings',
        name: 'Settings',
        component: Settings,
        meta: {auth: 'yes'}
    },
]


const router = createRouter({
    history: createWebHistory(''),
    routes
})

router.beforeEach(async (to, from, next) => {
    if (to.meta.auth) {
        let authorized = false
        await AuthAdapter.check()
            .then((result) => authorized = result)
            .catch(err =>{
            })
        const auth = to.meta.auth
        if (auth === 'yes' && !authorized) { // не даем войти неавторизированным на страницы для авторизированных
            return next({name: 'Login'})
        }
        if(auth === 'no' && authorized){ // не даем войти авторизированному на страницы для неавторизированных
            return next({name: 'Index'})
        }
    }
    return next()
})

export default router
