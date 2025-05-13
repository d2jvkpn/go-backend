<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import Home from './pages/Home.vue'

const route = useRoute()

const layoutComponent = computed(() => {
  const layout = route.meta.layout

  if (layout === 'home') {
    return Home
  }

  return 'div'
})
</script>


<template>
<component :is="layoutComponent">
  <router-view />
</component>
</template>


<style scoped>
</style>
