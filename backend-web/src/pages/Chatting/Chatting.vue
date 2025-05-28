<script setup>
import { ref, nextTick } from 'vue'
import { Operation, Position, ChatDotSquare, Plus } from '@element-plus/icons-vue'
// ChatLineSquare, Edit, Folder, Menu

const isCollapsed = ref(false)

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

const input = ref('')
const messages = ref([])
const chatLog = ref(null)

const send = () => {
  if (!input.value.trim()) return

  messages.value.push({ role: 'user', content: input.value })
  messages.value.push({ role: 'assistant', content: input.value })
  input.value = ''

  // 等 DOM 更新后滚动到底部
  nextTick(() => {
    if (chatLog.value) {
      // chatLog.value.scrollTop = chatLog.value.scrollHeight
      chatLog.value.scrollTo({ top: chatLog.value.scrollHeight, behavior: 'smooth' })
    }
  })
}
</script>

<template>
<div class="chat-layout">
  <div :class="['sidebar', { collapsed: isCollapsed }]">

    <div class="collapse-header" @click="toggleSidebar">
      <template v-if="!isCollapsed">
        <span class="collapse-text">Chatting List</span>
        <el-icon class="collapse-icon"><Operation /></el-icon>
      </template>

      <template v-else>
        <el-icon class="collapse-icon center-icon"><Operation /></el-icon>
      </template>
    </div>

    <el-menu class="sidebar-menu" default-active="1" :collapse="isCollapsed" >
      <el-menu-item index="1">
        <el-icon><ChatDotSquare /></el-icon>
        <template #title>GPT Chat</template>
      </el-menu-item>

      <el-menu-item index="2">
        <el-icon><ChatDotSquare /></el-icon>
        <template #title>A.I.D.A</template>
      </el-menu-item>

      <el-menu-item index="3">
        <el-icon><ChatDotSquare /></el-icon>
        <template #title>Writing Assistant</template>
      </el-menu-item>
    </el-menu>
  </div>

  <div class="chat-main">
    <div class="chat-log" ref="chatLog">
      <div :class="['chat-message', msg.role]" v-for="(msg, i) in messages" :key="i">
        <div class="bubble">
          <strong>{{ msg.role === 'user' ? 'You' : 'AI' }}:</strong> {{ msg.content }}
        </div>
      </div>
    </div>

    <div class="input-container">
      <textarea :rows="2" class="input-box" placeholder="Ask anything...(Press Ctrl+Enter to Send)"
        v-model="input" @keyup.enter.ctrl="send"
      ></textarea>

      <div class="input-actions">
        <el-icon class="input-icon" @click="send"><Position /></el-icon>
        <el-icon class="plus-icon" @click="toggleExtras" title="More"><Plus /></el-icon>
      </div>
    </div>
  </div>
</div>
</template>


<style scoped>
.chat-layout {
  display: flex;
  height: 100%;
}

.sidebar {
  width: 220px;
  transition: width 0.3s ease, background-color 0.3s ease;
  overflow: hidden;
  background: #fafafa;
  border-right: 1px solid #e4e4e4;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.sidebar.collapsed {
  width: 64px;
}
.sidebar.collapsed .collapse-text {
  opacity: 0;
}

.collapse-header {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  cursor: pointer;
  height: 48px;
  transition: all 0.2s;
  background: white;
}
/*
.collapse-header:hover {
  background-color: #f5f5f5;
}
*/
.collapse-text {
  font-weight: bold;
  font-size: 14px;
  flex: 1;
  opacity: 1;
  transition: opacity 0.2s ease;
}
.collapse-icon {
  font-size: 18px;
}
.center-icon {
  margin: 0 auto;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
  height: 100%;
  max-width: 50rem;
  margin: 0 auto;
}

.chat-log {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 12px;
  border-bottom: 1px solid #eee;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.chat-message .bubble {
  max-width: 70%;
  padding: 10px 14px;
  border-radius: 10px;
  line-height: 1.5;
  word-break: break-word;
}
.chat-message.user {
  display: flex;
  justify-content: flex-end;
}
.chat-message.user .bubble {
  background-color: #cce6ff;
  color: #003366;
  border-top-right-radius: 0;
}

.chat-message.assistant {
  display: flex;
  justify-content: flex-start;
}
.chat-message.assistant .bubble {
  background-color: #f5f5f5;
  color: #333;
  border-top-left-radius: 0;
}

.input-container {
  display: flex;
  gap: 8px;
  border: 1px solid #ccc;
  border-radius: 5px;
  justify-content: center; 
  height: 5rem;
  padding: 8px;
}

.input-box {
  flex: 1;
  font-size: 16px;
  height: 100%;
  border: none;
  outline: none;
  resize: none;
  background-color: transparent;
  width: 100%;
  font-size: 16px;
  line-height: 1.5;
}

.input-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  margin-left: 8px;
  padding-bottom: 10px;
}

.input-icon, .plus-icon {
  font-size: 20px;
  cursor: pointer;
  color: grey;
  transition: transform 0.2s;
}
.input-icon:hover, .plus-icon:hover {
  transform: scale(1.2);
}
</style>
