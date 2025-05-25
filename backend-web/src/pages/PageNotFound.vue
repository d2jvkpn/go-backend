<script setup lang="js">
import { useRoute } from 'vue-router'

const route = useRoute()
const fullPath = route.fullPath;
</script>

<template>
<div class="flex flex-col items-center justify-center gap-8 min-h-screen">
  <p> Page not found! You tried to visit: <code>{{ fullPath }}</code> </p>
  <el-button type="secondary" @click="$router.push('/')">Go Home</el-button>
</div>
</template>

<!-- bg-gradient-to-br from-green-500 to-sky-400 -->

<style scoped>
</style>
