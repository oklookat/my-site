import { createApp } from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import ElvenProgress from './common/plugins/ElvenProgress/ElvenProgress.js'

const app = createApp(App).use(store).use(router).use(ElvenProgress, {progressBarColor: 'white'}).mount('#app')
export default app
