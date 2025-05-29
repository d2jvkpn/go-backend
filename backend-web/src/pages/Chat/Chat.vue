<script setup>
import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Operation, Position, ChatSquare, DocumentAdd, Plus, MoreFilled } from '@element-plus/icons-vue'
import { Folder, Microphone, VideoPause, Delete } from '@element-plus/icons-vue'
// ChatLineSquare, Edit, Folder, Menu

//
const isCollapsed = ref(false)

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

//
const input = ref('')
const messages = ref([])
const loading = ref(false)
const chatLog = ref(null)
const editingSessionId = ref(null)

const chatSessions = ref([
  { id: "sess-1", title: "GPT Chat" },
  { id: "sess-2", title: "A.I.D.A" },
  { id: "sess-3", title: "Writting Assistant" },
])

//
const selectedSessionId = ref(null)

function selectedSession (session) {
  selectedSessionId.value = session.id
}

async function handleCommand ({ action, session }) {
  switch (action) {
  case "rename":
    editingSessionId.value = session.id
    await nextTick()
    const inputEl = document.getElementById(`session-input-${session.id}`)
    inputEl?.focus()
    break
  case "delete":
    chatSessions.value = chatSessions.value.filter(v => v.id !== session.id)
    break
  }
}

//
function addNewChat() {
  const newId = 'sess_' + Date.now()
  const newSession = { id: newId, title: '', editing: true }
  chatSessions.value.unshift(newSession)
  editingSessionId.value = newId
  nextTick(() => {
    const inputEl = document.getElementById(`session-input-${newId}`)
    inputEl?.focus()
  })
}

function finishEditing(session) {
  if (!session.title.trim()) {
    session.title = 'Untitled'
  }
  editingSessionId.value = null
}

//
async function cancelRequest() {
  console.log("==> cancelRequest")
  loading.value = false

  const lastMsg = messages.value[messages.value.length - 1]
  if (lastMsg?.role === 'assistant') {
     lastMsg.content = '<Canceled>'
  }
}

function scrollChat() {
  nextTick(() => {
    if (chatLog.value) {
      // chatLog.value.scrollTop = chatLog.value.scrollHeight
      chatLog.value.scrollTo({ top: chatLog.value.scrollHeight, behavior: 'smooth' })
    }
  })
}

function deleteMessage(msg, i) {
  if (msg.content == '__TYPING_DOTS__') {
    return;
  }

  ElMessage.warning(`removed message: ${msg.content.slice(0, 10)}`);
  messages.value.splice(i, 1)
}

async function sendRequest () {
  console.log("==> sendRequest")
  if (!input.value.trim()) {
    ElMessage.warning('Please enter a message before sending.')
    return
  }

  const msg = input.value;
  const ans = { role: 'assistant', content: '__TYPING_DOTS__' };
  input.value = ''
  messages.value.push({ role: 'user', content: msg }, ans);
  loading.value = true
  scrollChat()

  setTimeout(() => {
    // TODO: axios request
    ans.content = msg
    loading.value = false
    nextTick(() => scrollChat())
  }, 3000);
}

//
async function addDocs() {
  ElMessage.warning("TODO: add documents")
}

function listDocs() {
  ElMessage.warning("TODO: list documents")
}

/*
let controller = null;

function send() {
  // 创建新的 AbortController
  controller = new AbortController()

  try {
    const response = await axios.post('/api/ai-reply', { prompt: msg }, {
      signal: controller.signal,
    })

    ans.content = response.data.reply
  } catch (err) {
    if (axios.isCancel(err) || err.name === 'CanceledError' || err.code === 'ERR_CANCELED') {
      ans.content = '（已取消）'
    } else {
      ans.content = '请求失败，请稍后重试。'
      console.error(err)
    }
  } finally {
    loading.value = false
    controller = null
    scrollChat()
  }
}

function cancelRequest() {
  if (controller) {
    controller.abort()
    loading.value = false
  }
}
*/
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

    <div class="sidebar-menu" default-active="1" :collapse="isCollapsed" >
      <div index="0" class="new-chat" @click="addNewChat">
        <el-icon><Plus /></el-icon> New Chat
      </div>

      <div class="chat-session" :class="{ selected: selectedSessionId === session.id }" 
        :title="session.title || 'Untitled'"
        v-for="(session, index) in chatSessions" :index="index+1" :key="session.id"
      >
        <!--el-icon> <ChatSquare /> </el-icon-->

        <input
          placeholder="Enter title" style="width: 100%; border: none; outline: none; font-size: 12px"
          v-model="session.title" :id="`session-input-${session.id}`" v-if="editingSessionId === session.id"
          @blur="finishEditing(session)" @keyup.enter="finishEditing(session)"
        />

        <span class="session-title" v-if="editingSessionId !== session.id" @click="selectedSession(session)">
          {{ session.title || 'Untitled' }}
        </span>

        <el-dropdown @command="handleCommand" trigger="click">
          <el-icon style="margin-left: auto; cursor: pointer;"> <MoreFilled /> </el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="{ action: 'rename', session }"> Rename </el-dropdown-item>
              <el-dropdown-item :command="{ action: 'delete', session }"> Delete </el-dropdown-item>
              </el-dropdown-menu>
          </template>
       </el-dropdown>

      </div>
    </div>
  </div>

  <div class="chat-main">

    <div class="chat-log" ref="chatLog">
      <div :class="['chat-message', msg.role]" v-for="(msg, i) in messages" :key="i">
        <div class="bubble g-cell-center">
          <el-icon class="delete-icon" @click="deleteMessage(msg, i)"> <Delete /> </el-icon>
          <template v-if="msg.content === '__TYPING_DOTS__'"> <TypingDots /></template>
          <template v-else> {{ msg.content }} </template>
        </div>
      </div>
    </div>

    <div class="input-container" title="Press Ctrl+Enter to Send">
      <textarea class="input-box" rows="2"
        :placeholder="loading ? 'Thinking...' : 'Ask anything...'"
        v-model="input" @keyup.enter.ctrl="sendRequest" :disabled="loading"
      ></textarea>

      <div class="input-actions">
        <el-icon @click="" title="More"><Microphone /></el-icon>
        <el-icon  @click="loading ? cancelRequest() : sendRequest()">
          <template v-if="loading"> <VideoPause /> </template>
          <template v-else> <Position /> </template>
        </el-icon>
        <el-icon @click="listDocs" title="More"><Folder /></el-icon>
        <el-icon @click="addDocs" title="More"><DocumentAdd /></el-icon>
      </div>
    </div>

  </div>
</div>
</template>


<style scoped>
.chat-layout {
  display: flex;
  /* height: 100%; */
  height: 100vh;
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

.new-chat {
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.2s;
}

.new-chat:hover {
  background: #cce6ff;
}

.chat-session {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  margin: 0.21rem;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.2s;
}

.chat-session:hover {
  background-color: #e0e0e0;
}

.chat-session.selected {
  background-color: #cce6ff;
}

.session-title {
  font-size: 14px;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  cursor: pointer;
  user-select: none;
  transition: color 0.2s;
}

.bubble {
  position: relative;
  padding-right: 1.5em; /* 给右上角留空间 */
}

.delete-icon {
  position: absolute;
  top: 6px;
  right: 6px;
  cursor: pointer;
  display: none;
  color: #888;
  font-size: 16px;
}

.bubble:hover .delete-icon {
  display: inline-block;
}

.chat-message .bubble {
  max-width: 70%;
  padding: 10px 14px;
  border-radius: 10px;
  line-height: 1.5;
  word-break: break-word;
  white-space: pre-wrap;
}
.chat-message.user {
  display: flex;
  justify-content: flex-end;
}
.chat-message.user .bubble {
  background-color: #f5f5f5;
  color: #333;
  border-top-right-radius: 0;
}

/* an element has both class chat-message and assistant */
.chat-message.assistant {
  display: flex;
  justify-content: flex-start;
}
.chat-message.assistant .bubble {
  background-color: #cce6ff;
  color: #003366;
  border-top-left-radius: 0;
}

.input-container {
  display: flex;
  gap: 8px;
  border: 1px solid #ccc;
  border-radius: 10px;
  justify-content: center; 
  height: 6rem;
  padding: 10px;
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
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin: auto 2px auto 5px;
}

/* all direct children of class input-actions */
.input-actions > * {
  font-size: 24px;
  cursor: pointer;
  color: grey;
  transition: transform 0.2s;
}
.input-actions > *:hover {
  transform: scale(1.2);
}
</style>
