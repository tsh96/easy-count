import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { routes } from 'vue-router/auto-routes'
import { createRouter, createWebHashHistory } from 'vue-router'
import { VueQrcodeReader } from 'vue-qrcode-reader'
import { isAuthenticated } from './composables/auth'

const router = createRouter({
  history: createWebHashHistory(),
  // pass the generated routes written by the plugin 🤖
  routes,
})

// Navigation guard for authentication
router.beforeEach((to, _from, next) => {
  // Public routes that don't require authentication
  const publicPages = ['/auth'];
  const authRequired = !publicPages.includes(to.path);

  if (authRequired && !isAuthenticated.value) {
    // Redirect to login if not authenticated
    next('/auth');
  } else if (to.path === '/auth' && isAuthenticated.value) {
    // Redirect to home if already authenticated
    next('/');
  } else {
    next();
  }
});

const meta = document.createElement('meta')
meta.name = 'naive-ui-style'
document.head.appendChild(meta)

createApp(App).use(router).use(VueQrcodeReader).mount('#app')
