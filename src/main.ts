import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { routes } from 'vue-router/auto-routes'
import { createRouter, createWebHashHistory } from 'vue-router'
import { VueQrcodeReader } from 'vue-qrcode-reader'

const router = createRouter({
  history: createWebHashHistory(),
  // pass the generated routes written by the plugin 🤖
  routes,
})

const meta = document.createElement('meta')
meta.name = 'naive-ui-style'
document.head.appendChild(meta)

createApp(App).use(router).use(VueQrcodeReader).mount('#app')
