import { createStore } from 'vuex'
import createPersistedState from "vuex-persistedstate"
import {authStore} from "@/store/modules";

const authState = createPersistedState({
  key: 'elven',
  paths: ['authStore']
})

export default createStore({
  state: {
  },
  mutations: {
  },
  actions: {
  },
  modules: {
    authStore,
  },
  plugins: [authState],
})
