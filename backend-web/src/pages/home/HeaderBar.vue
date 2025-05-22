<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { ArrowDown, Fold, Expand } from '@element-plus/icons-vue'

import ChangePassword from './ChangePassword.vue'
// import { useUserStore } from '@/js/stores/user'
import { logout } from "@/js/_login.js"
import { clearAccount } from "@/js/stores/local.js"
import { service } from  "@/js/utils/request.js"

//
const props = defineProps({
  firstname: String,
  lastname: String,
  isSidebarHidden: Boolean,
})

const showChangePassword = ref(false)

const handlePasswordSubmit = () => {
  showChangePassword.value = false
}

//
const emit = defineEmits(['toggleSidebar'])

//
const router = useRouter()
// const userStore = useUserStore()

const confirmLogoutV1 = () => {
  ElMessageBox.confirm(
    'Are you sure you want to log out?',
    'Logout Confirmation',
    { confirmButtonText: 'Logout', cancelButtonText: 'Cancel', type: 'warning' }
  )
  .then(() => {
    // userStore.logout()
    // localStorage.removeItem('token')
    clearAccount()
    router.push('/login')
    ElMessage.success('You have been logged out.')
  })
  .catch(() => {
    ElMessage.info('Logout canceled.')
  })
}

const confirmLogout = async () => {
  try {
    await ElMessageBox.confirm(
      'Are you sure you want to log out?',
      'Logout Confirmation',
      { confirmButtonText: 'Logout', cancelButtonText: 'Cancel', type: 'warning' },
    );

    await service.post(`${import.meta.env.VITE_API_URL}/api/v1/auth/account/logout`)

    clearAccount()
    router.push('/login')
    ElMessage.success('You have been logged out.')
  } catch (err) {
    ElMessage.error(err.response?.data?.msg || 'Logout failed')
  } finally {
    // TODO
  }
}

const handleCommand = (command) => {
  switch (command) {
  case 'profile':
    router.push('/home/settings/profile')
    break
  case 'change_password':
    showChangePassword.value = true
    break
  case 'logout':
    confirmLogout()
    break
  default:
    alert(`!!! unknown command: ${command}`)
  }
}

</script>


<template>
<header class=headerbar>
  <div class="headerbar-left">
    <el-button text circle @click="$emit('toggleSidebar')" class="headerbar-toggle-btn">
      <el-icon> <component :is="isSidebarHidden ? Expand : Fold" /> </el-icon>
    </el-button>
    <div class="headerbar-logo"> 🌀 Home </div>
  </div>

  <el-dropdown @command="handleCommand">
    <span class="headerbar-account el-dropdown-link">
      {{ firstname }} {{ lastname }} <el-icon> <ArrowDown /> </el-icon>
    </span>

    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="profile"> 👤 Profile </el-dropdown-item>
        <el-dropdown-item command="change_password"> 🔒 Change password </el-dropdown-item>
        <el-dropdown-item divided command="logout"> ⏻ Logout </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>

<ChangePassword :visible="showChangePassword" @close="showChangePassword = false"/>
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
