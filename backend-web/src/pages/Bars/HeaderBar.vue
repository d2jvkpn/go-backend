<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { ArrowDown, Fold, Expand , ChatLineSquare, Postcard, Lock, SwitchButton } from '@element-plus/icons-vue'

import ChangePassword from './ChangePassword.vue'
import stores from "@/js/stores"
import { request } from  "@/js/api"

//
const props = defineProps({
  account: { type: Object },
  isSidebarHidden: { type: Boolean },
})

const emit = defineEmits(['toggleSidebar'])

const router = useRouter()
const showChangePassword = ref(false)

async function confirmLogout() {
  let LoggedOut = true;

  try {
    await ElMessageBox.confirm(
      'Are you sure you want to log out?',
      'Logout Confirmation',
      { confirmButtonText: 'Logout', cancelButtonText: 'Cancel', type: 'warning' },
    );

    await request.post(`${import.meta.env.VITE_API_URL}/api/v1/auth/account/logout`);
    ElMessage.success('You have been logged out.');
  } catch (err) {
    if (err === "cancel") {
      LoggedOut = false;
      ElMessage.info('Logout canceled');
    } else {
      ElMessage.error(err.response?.data?.msg || 'Logout failed');
    }
  } finally {
  }

  if (LoggedOut) {
    // Clear the token first to avoid automatic redirection to /dashboard and unintended API requests after landing on the /login page.
    stores.clearAccount();
    router.push('/login');
  }
}

function handleCommand (command) {
  switch (command) {
  case "chat":
    window.open(`${import.meta.env.VITE_BASE_PATH}/chat`, '_blank')
    break
  case 'profile':
    router.push('/dashboard/settings/profile')
    break
  case 'change_password':
    showChangePassword.value = true
    break
  case 'logout':
    confirmLogout()
    break
  default:
    console.error(`!!! unknown command: ${command}`)
  }
}

</script>


<template>
<header class=headerbar>
  <div class="headerbar-left">
    <el-button text circle @click="$emit('toggleSidebar')" class="headerbar-toggle-btn">
      <el-icon> <component :is="isSidebarHidden ? Expand : Fold" /> </el-icon>
    </el-button>
    <div class="headerbar-logo"> 🌀 Dashboard </div>
  </div>

  <el-dropdown @command="handleCommand">
    <span class="headerbar-account el-dropdown-link">
      {{ account.firstname }} {{ account.lastname }} <el-icon> <ArrowDown /> </el-icon>
    </span>

    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="chat">
          <el-icon> <ChatLineSquare /> </el-icon> Chat
        </el-dropdown-item>
        <!--el-dropdown-item command="profile"> 👤 Profile </el-dropdown-item-->
        <!--el-dropdown-item command="change_password"> 🔒 Change password </el-dropdown-item-->
        <!--el-dropdown-item divided command="logout"> ⏻ Logout </el-dropdown-item-->

        <el-dropdown-item command="profile">
          <el-icon> <Postcard /> </el-icon> Profile
        </el-dropdown-item>

        <el-dropdown-item command="change_password">
          <el-icon><Lock /></el-icon> Change password
        </el-dropdown-item>

        <el-dropdown-item divided command="logout">
          <el-icon><SwitchButton /></el-icon> Logout
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>

  <ChangePassword v-model:visible="showChangePassword" />
</header>
</template>


<style scoped>
.headerbar {
  background-color: #fff;
  padding: 0 0.8rem;
  font-size: 2rem;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.headerbar-left {
  display: flex;
  align-items: center;
  gap: 0.2rem;
}

.headerbar-toggle-btn {
  font-size: 1.4rem;
  /*
  position: absolute;
  top: 2px;
  left: 2px;
  z-index: 10;
  */
}

.headerbar-logo {
  font-size: 1.0rem;
  color: #409eff;
  font-weight: bold;
}

.headerbar-account {
  cursor: pointer;
  color: #333;
  display: flex;
  align-items: center;
  gap: 4px;
}

.headerbar-account.el-dropdown-link {
  transition: font-size 0.2s ease;
  border: none !important;
}

.headerbar-account.el-dropdown-link:hover {
  border: none !important;
  background: none !important;
  font-size: 16px;
  color: #409eff;
}
</style>
