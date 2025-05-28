<template>
  <span class="typing-dots">
    <span v-for="i in 10" :key="i">.</span>
  </span>
</template>

<script setup>
// 无需逻辑
</script>

<style scoped>
.typing-dots {
  display: inline-block;
  font-size: 16px;
  color: #999;
}

.typing-dots span {
  opacity: 0;
  animation: blink 1.5s infinite;
  animation-fill-mode: both;
  display: inline-block;
}

.typing-dots span:nth-child(1) { animation-delay: 0s; }
.typing-dots span:nth-child(2) { animation-delay: 0.2s; }
.typing-dots span:nth-child(3) { animation-delay: 0.4s; }
.typing-dots span:nth-child(4) { animation-delay: 0.6s; }
.typing-dots span:nth-child(5) { animation-delay: 0.8s; }
.typing-dots span:nth-child(6) { animation-delay: 1s; }
.typing-dots span:nth-child(7) { animation-delay: 1.2s; }
.typing-dots span:nth-child(8) { animation-delay: 1.4s; }
.typing-dots span:nth-child(9) { animation-delay: 1.6s; }
.typing-dots span:nth-child(10) { animation-delay: 1.8s; }

@keyframes blink {
  0%, 100% { opacity: 0; }
  50%      { opacity: 1; }
}
</style>
