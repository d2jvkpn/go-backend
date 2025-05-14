<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { ArrowDown, Fold, Expand } from '@element-plus/icons-vue'

import ChangePassword from './ChangePassword.vue'
// import { useUserStore } from '@/js/stores/user'

//
const props = defineProps({
  accountName: String,
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

const confirmLogout = () => {
  ElMessageBox.confirm(
    'Are you sure you want to log out?',
    'Logout Confirmation',
    {
      confirmButtonText: 'Logout',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  )
  .then(() => {
    // userStore.logout()
    localStorage.clear()
    router.push('/login')
    ElMessage.success('You have been logged out.')
  })
  .catch(() => {
    ElMessage.info('Logout canceled.')
  })
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
      // localStorage.removeItem('token')
      // localStorage.clear()
      // router.push('/login')
      confirmLogout()
      break
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
      {{ accountName }} <el-icon> <ArrowDown /> </el-icon>
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
