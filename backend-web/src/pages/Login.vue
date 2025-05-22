<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { login } from "@/js/_login.js"
import { service } from  "@/js/utils/request.js"
import { setAccount } from "@/js/stores/local.js"

// console.log(`==> import.meta.env: ${JSON.stringify(import.meta.env)}`);

const account = ref('')
const password = ref('')
const loading = ref(false)
const router = useRouter()

const sumbitLoginV1 = () => {
  if (!account.value || !password.value) {
    alert('Please enter acocunt and password!')
    return
  }

  /*
  let firstname = "Jane";
  let lastname = "Doe";
  localStorage.setItem('fristname', "Jane")
  localStorage.setItem('lastname', "Doe")
  localStorage.setItem('token', `${account.value}:${password.value}`)

  localStorage.setItem('level', "admin")

  ElMessage.success(`Welcome back, ${firstname} ${lastname}!`)
  router.push('/home/dashboard')
  */

  const callback = (data) => {
    localStorage.setItem('firstname', data.firstname)
    localStorage.setItem('lastname', data.lastname)
    localStorage.setItem('token', data.token)
    localStorage.setItem('level', data.level)

    loading.value = false
    ElMessage.success(`Welcome back, ${data.firstname} ${data.lastname}!`)
    router.push('/home/dashboard')
  }

  const onError = () => {
    loading.value = false
  }

  const data = { password: password.value }
  if (account.value.includes("@")) {
    data.email = account.value;
  } else {
    data.phone = account.value;
  }

  login(data, callback, onError)
}

const sumbitLogin = async () => {
  if (!account.value || !password.value) {
    alert('Please enter acocunt and password!')
    return
  }

  const loginData = { password: password.value }
  if (account.value.includes("@")) {
    loginData.email = account.value;
  } else {
    loginData.phone = account.value;
  }

  try {
    const data = await service.post(
      `${import.meta.env.VITE_API_URL}/api/v1/open/account/login?platform=web`,
      loginData,
    )

    setAccount(data)
    ElMessage.success(`Welcome back, ${data.firstname} ${data.lastname}!`)
    router.push('/home/dashboard')
  } catch (err) {
    console.log(`!!! login error: ${JSON.stringify(err)}`)
  } finally {
    loading.value = false
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
        <el-input type="password" v-model="password" placeholder="password" show-password />
      </el-form-item>

      <div class="login-button">
        <el-button type="primary" @click="sumbitLogin"> {{ loading ? 'Logining...' : 'Login' }} </el-button>
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
