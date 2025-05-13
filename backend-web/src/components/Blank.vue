<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()
console.log(`==> Page: title=${route.meta.title}, roles=${JSON.stringify(route.meta.roles)}`)

const props = defineProps({
  text: { type: String, default: 'Blank' },
  color: { type: String, default: 'green' }
})
</script>


<template>
  <div :style="{ color: color }"> Welcome to {{ text }}! </div>
</template>

<style scoped>
/*
div {
  color: green;
}
*/
</style>
