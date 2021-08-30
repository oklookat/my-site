import { createApp } from 'vue'
import App from '@/App.vue'
import router from "@/router"
import store from '@/store'
import ElvenProgress from '@/common/plugins/ElvenProgress/ElvenProgress.js'
import ElvenPlayer from '@/common/plugins/ElvenPlayer/ElvenPlayer.js'
import ElvenNotify from '@/common/plugins/ElvenNotify/ElvenNotify.js'

const app = createApp(App)
    .use(store)
    .use(router)
    .use(ElvenProgress)
    .use(ElvenPlayer)
    .use(ElvenNotify)
    .mount('#app')

window.app = app
