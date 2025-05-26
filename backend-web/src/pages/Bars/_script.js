import { useUserStore } from '@/js/stores/user'
import { logout } from "@/js/api/_login.js"

// const userStore = useUserStore()

const confirmLogout = () => {
  ElMessageBox.confirm(
    'Are you sure you want to log out?',
    'Logout Confirmation',
    { confirmButtonText: 'Logout', cancelButtonText: 'Cancel', type: 'warning' },
  )
  .then(() => {
    // userStore.logout()
    // localStorage.removeItem('token')
    stores.clearAccount();
    router.push('/login');
    ElMessage.success('You have been logged out.');
  })
  .catch(() => {
    ElMessage.info('Logout canceled.');
  })
}
