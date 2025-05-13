#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _dir=$(readlink -f `dirname "$0"`)


function read_embeded() {
    key=$1
    sed -n "/^__START_${key}__$/,/^__END_${key}__/p" "$0" | tail -n +2 | head -n -1
}

mkdir -p archive/src/components src/router src/styles src/layout src/pages src/utils
mv vite.config.js archive/
mv src/App.vue archive/src/
mv src/style.css archive/src/
mv src/main.js archive/src/
mv src/components/HelloWorld.vue archive/src/components/


read_embeded Vite > vite.config.js
read_embeded App > src/App.vue
read_embeded Main > src/main.js
read_embeded Style > src/styles/style.css
read_embeded Hello > src/pages/Hello.vue
read_embeded Router > src/router/index.js


exit 0

# vite.config.js
__START_Vite__
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

import vueDevTools from 'vite-plugin-vue-devtools'
import { fileURLToPath, URL } from 'node:url'
import tailwindcss from "@tailwindcss/vite"
import Components from 'unplugin-vue-components/vite'
//import AutoImport from 'unplugin-auto-import/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'


// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),

    vueDevTools(),
    tailwindcss(),
    //AutoImport({resolvers: [ElementPlusResolver()]}),
    Components({resolvers: [ElementPlusResolver()]}),
  ],

  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      // extensions: ['.mjs', '.js', '.jsx', '.json', '.vue'],
    },
  },
})
__END_Vite__


__START_App__
<template>
<div id="app">
  <!--
  <nav>
    <router-link to="/"> Login  </router-link>
    <router-link to="/about"> Abount </router-link>
  </nav>
  -->

  <router-view />
</div>
</template>

<!--style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style-->
__END_App__


__START_Hello__
<script setup lang="js">
</script>

<template>
<div
  class="flex flex-col items-center justify-center gap-8 min-h-screen bg-gradient-to-br from-green-500 to-sky-400"
>
  <p>Hello, world!</p>
</div>
</template>

<style scoped>
</style>
__END_Hello__


__START_Router__
import { createRouter, createWebHistory } from 'vue-router'
import Hello from '../pages/Hello.vue'

const routes = [
  {
    path: '/hello', name: 'Hello', component: Hello,
    meta: { title: 'Vue - Hello' },
  },
  {
    path: '/world', name: 'World',
    component: () => import('../pages/Hello.vue'),
    meta: { title: 'Vue - World' },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/hello',
  },
  /*
  {
    path: '/:pathMatch(.*)*', name: 'NotFound',
    component: () => import('../pages/NotFound.vue'),
  }
  */
]

const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_BASE_PATH),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = String(to.meta.title) || 'Page Not Found';
  next();
})

export default router;
__END_Router__


__START_Main__
import { createApp } from 'vue'
//import './style.css'
import App from './App.vue'

import ElementPlus from 'element-plus'
import './styles/style.css'
import router from './router'


// createApp(App).mount('#app')

const app = createApp(App)

app.use(ElementPlus)
app.use(router)

app.mount('#app')
__END_Main__


__START_Style__
@import "tailwindcss";
/*
@tailwind base;
@tailwind components;
@tailwind utilities;
*/

@import "element-plus/dist/index.css";
__END_Style__
