<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { login } from "../js/login.js"

// console.log(`==> import.meta.env: ${JSON.stringify(import.meta.env)}`);

const account = ref('')
const password = ref('')
const loading = ref(false)
const router = useRouter()

const sumbitLogin = () => {
  if (account.value && password.value) {
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
      localStorage.setItem('fristname', data.firstname)
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

    login({email: account.value, password: password.value}, callback, onError)
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
