import { createRouter, createWebHistory } from "vue-router"
import Login from "@/pages/Login.vue"
import PageNotFound from "@/pages/PageNotFound.vue"

export const allRoutes = [
  { path: "/", redirect: "/login" },

  {
    path: "/login", name: "Login",
    meta: { title: "Backend - Login", requiresAuth: false, levels: ["any"] },
    component: Login,
  },

  {
    path: "/home/dashboard", name: "Dashboard",
    meta: { title: "Backend - Dashboard", requiresAuth: true, levels: ["editor", "admin"] },
    component: () => import("@/pages/Dashboard/Dashboard.vue"),
  },

  {
    path: "/home/accounts", name: "Accounts",
    meta: { title: "Backend - Accounts", requiresAuth: true, levels: ["admin"] },
    component: () => import("@/pages/Accounts/Accounts.vue"),
  },

  {
    path: "/home/settings", name: "Settings",
    meta: { title: "Backend - Settings", requiresAuth: true, levels: ["any"] },
    children: [
      {
        path: "profile", name: "Profile",
        meta: { title: "Backend - Profile", requiresAuth: true, levels: ["any"] },
        component: () => import("@/pages/Profile/Profile.vue"),
      },
      {
        path: "security", name: "Security",
        meta: { title: "Backend - Security", requiresAuth: true, levels: ["admin"] },
        component: () => import("@/pages/Security/Security.vue"),
      },
    ],
  },

  {
    path: "/page-not-found", name: "PageNotFound",
    meta: { title: "Backend - Page not found", requiresAuth: false, levels: ["any"] },
    component: PageNotFound,
  },

  { path: "/:pathMatch(.*)*", redirect: "/page-not-found" },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_BASE_PATH),
  routes: allRoutes,
})

router.beforeEach((to, from, next) => {
  document.title = String(to.meta.title) || "Page Not Found";

  const isLoggedIn = !!localStorage.getItem("token");

  if (to.meta.requiresAuth && !isLoggedIn) {
    return next("/login")
  }

  /* TODO: pemission denied
  let level = new Set(JSON.parse(localStorage.getItem("level")));
  if (item.meta.levels?.includes("any") || item.meta.levels?.some(v => props.levels.has(v))) {
    return next("/pemission-denied")
  }
  */

  next()
})

export default router
