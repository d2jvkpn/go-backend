<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

// console.log(`==> import.meta.env: ${JSON.stringify(import.meta.env)}`);

const account = ref('')
const password = ref('')
const router = useRouter()

const login = () => {
  if (account.value && password.value) {
    let username = "Jane Doe"
    localStorage.setItem('token', `${account.value}:${password.value}`)
    localStorage.setItem('accountName', username)
    localStorage.setItem('roles', JSON.stringify(["admin"]))

    ElMessage.success(`Welcome back, ${username}!`)
    router.push('/home/dashboard')
  } else {
    alert('Please enter acocunt and password!')
  }
}
</script>


<template>
<div class="login">
  <el-card class="login-card">
    <template #header> Please Login </template>
    <el-form @submit.prevent="login">
      <el-form-item>
        <el-input v-model="account" placeholder="email or phone" />
      </el-form-item>

      <el-form-item>
        <el-input type="password" v-model="password" placeholder="password" />
      </el-form-item>

      <div class="login-button">
        <el-button type="primary" @click="login"> Login </el-button>
      </div>
    </el-form>
  </el-card>
</div>
</template>


<style scoped>
.login {
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.login-card {
  width: 30rem;
  margin: 0rem 2rem;
}

.login-button {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
