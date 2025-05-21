<script setup>
import { ref, onMounted, watchEffect} from 'vue'
import { Fold, Expand, ArrowDown } from '@element-plus/icons-vue'

import Sidebar from "./home/Sidebar.vue"
import HeaderBar from "./home/HeaderBar.vue"

//
const accountName = localStorage.getItem('accountName')
const level = localStorage.getItem('level') // JSON.parse(...)

//
const isHidden = ref(false)

onMounted(() => {
  const mediaQuery = window.matchMedia('(max-width: 720px)')

  isHidden.value = mediaQuery.matches

  mediaQuery.addEventListener('change', (e) => {
    isHidden.value = e.matches
  })
})

const toggleCollapse = () => {
  isHidden.value = !isHidden.value
}
</script>


<template>
<div class="home">
  <HeaderBar
    class="home-headerbar"
    :accountName="accountName"
    :isSidebarHidden="isHidden"
    @toggleSidebar="isHidden = !isHidden"
  />

  <div class="home-main">
    <Sidebar class="home-sidebar" v-show="!isHidden" :level="level" />

    <main class="home-content">
      <router-view v-slot="{ Component }">
        <KeepAlive> <component :is="Component" /> </KeepAlive>
      </router-view>
    </main>
  </div>
</div>
</template>


<style scoped>
.home {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.home-headerbar {
  height: 60px;
}

.home-main {
  display: flex;
  flex: 1; /* flex-grow: 1; flex-shrink: 1; flex-basis: 0%; */
  overflow: hidden;
}

.home-sidebar {
  width: 15rem;
}

.home-content {
  position: relative;
  flex: 1;
  padding: 20px;
  background-color: #f5f7fa;
  overflow-y: auto;
}
</style>
