// $ npm install pinia pinia-plugin-persistedstate


// path: src/js/stores/user.ts
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  persist: true
  state: () => ({
    userId: null,
    username: '',
    roles: [],
    token: '',
  }),
  actions: {
    login(payload) {
      this.userId = payload.id
      this.username = payload.username
      this.roles = payload.roles
      this.token = payload.token
    },
    logout() {
      this.$reset()
    },
  },
})


// path: src/main.ts
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

app.use(pinia)

// application
import { useUserStore } from '@/stores/user.js'

const userStore = useUserStore()

console.log(userStore.roles)


// path: src/router/index.js
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.meta.roles && !to.meta.roles.some(role => userStore.roles.includes(role))) {
    return next('/403')
  }
  next()
})
