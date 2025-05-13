import { createApp } from 'vue'
//import './style.css'
import App from './App.vue'

import './styles/style.css'
import router from './router'

import ElementPlus from 'element-plus'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

//createApp(App).mount('#app')

const app = createApp(App)

app.use(ElementPlus)

app.use(router)

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
app.use(pinia)

app.mount('#app')
