<script setup>
import { ref, onMounted, watchEffect} from 'vue'
import { Fold, Expand, ArrowDown } from '@element-plus/icons-vue'

import Sidebar from "./Bars/Sidebar.vue"
import HeaderBar from "./Bars/HeaderBar.vue"
import { getAccount } from "@/js/stores/storage.js"

//
const account = getAccount();

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
<div class="dashboard">
  <HeaderBar
    class="dashboard-headerbar"
    :account="account"
    :isSidebarHidden="isHidden"
    @toggleSidebar="isHidden = !isHidden"
  />

  <div class="dashboard-main">
    <Sidebar class="dashboard-sidebar" v-show="!isHidden" :account="account" />

    <main class="dashboard-content">
      <router-view v-slot="{ Component }">
        <KeepAlive> <component :is="Component" /> </KeepAlive>
      </router-view>
    </main>
  </div>
</div>
</template>


<style scoped>
.dashboard {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.dashboard-headerbar {
  height: 60px;
}

.dashboard-main {
  display: flex;
  flex: 1; /* flex-grow: 1; flex-shrink: 1; flex-basis: 0%; */
  overflow: hidden;
}

.dashboard-sidebar {
  width: 15rem;
}

.dashboard-content {
  position: relative;
  flex: 1;
  padding: 20px;
  background-color: #f5f7fa;
  overflow-y: auto;
}
</style>
