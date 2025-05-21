<script setup>
import { computed, ref, onBeforeMount } from 'vue'
import { useRoute } from 'vue-router'

import Home from './pages/Home.vue'

//
console.log(`==> import.meta.env: ${JSON.stringify(import.meta.env)}`);

const config = ref({});

onBeforeMount(() => {
  fetch(`${import.meta.env.VITE_BASE_PATH}/app.json`)
    .then((response) => response.json())
    .then((data) => {
      config.value = data;
      console.log(`==> Got app.json: ${JSON.stringify(data)}`);
    })
    .catch(error => console.error(`!!! Error loading app.json: ${error}`));
});

//
const route = useRoute()

const layoutComponent = computed(() => {
  const requiresAuth = route.meta.requiresAuth

  if (requiresAuth) {
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
