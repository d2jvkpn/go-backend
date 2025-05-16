<script setup>
import { computed, ref, onBeforeMount } from 'vue'
import { useRoute } from 'vue-router'

import Home from './pages/Home.vue'

//
console.log(`==> import.meta.env: ${JSON.stringify(import.meta.env)}`);

const config = ref({});

onBeforeMount(() => {
  fetch('config.json')
    .then(response => response.json())
    .then((data) => {
      config.value = data;
      console.log(`==> Got config: ${JSON.stringify(data)}`);
    })
    .catch(error => console.error(`!!! Error loading config: ${error}`));
});

//
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
