import { createRouter, createWebHistory } from 'vue-router'
import PageNotFound from '@/pages/PageNotFound.vue'

export const allRoutes = [
  { path: '/', redirect: '/login' },

  {
    path: '/login',
    name: "Login",
    meta: { title: "App - Login", layout: 'none', roles: ["any"] },
    component: () => import('@/pages/Login.vue'),
  },

  {
    path: '/home/dashboard',
    name: "Dashboard",
    meta: { title: "App - Dashboard", layout: 'home', roles: ["editor", "admin"] },
    component: () => import('@/pages/home/Dashboard.vue'),
  },

  {
    path: '/home/accounts',
    name: "Accounts",
    meta: { title: "App - Accounts", layout: 'home', roles: ["admin"] },
    component: () => import('@/pages/home/Accounts.vue'),
  },

  {
    path: '/home/settings',
    name: "Settings",
    meta: { title: "App - Settings", layout: 'home', roles: ["any"] },
    children: [
      {
        path: 'profile',
        name: "Profile",
        meta: { title: "App - Profile", layout: 'home', roles: ["any"] },
        component: () => import('@/pages/home/Profile.vue'),
      },
      {
        path: 'security',
        name: "Security",
        meta: { title: "App - Security", layout: 'home', roles: ["admin"] },
        component: () => import('@/pages/home/Security.vue'),
      },
    ],
  },

  {
    path: '/page-not-found',
    name: "PageNotFound",
    meta: { title: "App - Page not found", layout: 'none', roles: ["any"] },
    component: () => PageNotFound,
  },

  { path: '/:pathMatch(.*)*', redirect: '/page-not-found' },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_BASE_PATH),
  routes: allRoutes,
})

router.beforeEach((to, from, next) => {
  document.title = String(to.meta.title) || 'Page Not Found';

  const isLoggedIn = !!localStorage.getItem('token');

  if (to.meta.layout === 'home' && !isLoggedIn) {
    return next('/login')
  }

  /* TODO: pemission denied
  let roles = new Set(JSON.parse(localStorage.getItem('roles')));
  if (!meta.roles?.includes("any") && !item.meta.roles?.some(e => roles.has(e))) {
    return next('/pemission-denied')
  }
  */

  next()
})

export default router
