import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/styles/main.css'
import { registerSW } from 'virtual:pwa-register'

// Register service worker — auto-updates in background
registerSW({ immediate: true })

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
