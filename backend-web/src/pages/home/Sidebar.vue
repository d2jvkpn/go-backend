<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Histogram, User, Setting, Postcard, Box } from '@element-plus/icons-vue'

import { allRoutes } from '@/router/index'

const route = useRoute()

const props = defineProps({
  level: String, // Set,
})

//console.log(`--> level: ${level}`)
//console.log(`~~~ ${allRoutes.length}`);

const visibleRouteNames = computed(() => {
  const result = []

  const collect = (items) => {
    for (const item of items) {
      if (!item.name || !item.meta) { // only named routers
        continue
      }

      // console.log(`~~~ ${item.name}`)
      if (item.meta.levels?.includes("any") || item.meta.levels?.includes(props.level)) {
        result.push(item.name)
      }

      if (item.children) {
        collect(item.children)
      }
    }
  }

  collect(allRoutes)
  // console.log(`--> ${JSON.stringify(result)}`)
  return result
})
</script>


<template>
<aside>
  <el-menu :default-active="$route.path" router>
    <el-menu-item index="/home/dashboard" v-if="visibleRouteNames.includes('Dashboard')">
      <el-icon> <Histogram /> </el-icon>
      Dashboard
    </el-menu-item>

    <el-menu-item index="/home/accounts" v-if="visibleRouteNames.includes('Accounts')">
      <el-icon> <User /> </el-icon>
      Accounts
    </el-menu-item>

    <el-sub-menu index="/home/settings" v-if="visibleRouteNames.includes('Settings')">
      <template #title>
        <el-icon> <Setting /> </el-icon>
        Settings
      </template>
      <el-menu-item index="/home/settings/profile" v-if="visibleRouteNames.includes('Profile')">
        <el-icon><Postcard /></el-icon>
        Profile
      </el-menu-item>

      <el-menu-item index="/home/settings/security" v-if="visibleRouteNames.includes('Security')">
        <el-icon><Box /></el-icon>
        Security
      </el-menu-item>
    </el-sub-menu>
  </el-menu>
</aside>
</template>


<style scoped>
aside {
  color: #001219;
  padding-top: 10px;
  overflow-y: auto;
}
</style>
