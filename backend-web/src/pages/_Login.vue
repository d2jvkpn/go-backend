<script setup>
import { ref } from 'vue'

const account = ref('')
const password = ref('')
const captcha = ref('')
const loading = ref(false)

const countdown = ref(0)
let timer = null

function sendCode() {
  // 模拟请求发送验证码
  console.log('Sending verification code to:', account.value)
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
    }
  }, 1000)
}

function submit() {
  loading.value = true
  console.log('Logging in with:', {
    account: account.value,
    password: password.value,
    captcha: captcha.value
  })

  setTimeout(() => {
    loading.value = false
  }, 1000)
}
</script>

<template>
  <div class="login">
    <el-card class="login-card">
      <template #header> Please Login </template>

      <el-form @submit.prevent="submit">
        <el-form-item>
          <el-input placeholder="email or phone" v-model="account" />
        </el-form-item>

        <el-form-item>
          <el-input
            placeholder="password"
            type="password"
            v-model="password"
            clearable
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-input
            v-model="captcha"
            placeholder="Verification Code"
            maxlength="6"
            style="width: 60%; margin-right: 10px"
          />
          <el-button type="primary" plain @click="sendCode" :disabled="countdown > 0">
            {{ countdown > 0 ? countdown + 's' : 'Send Code' }}
          </el-button>
        </el-form-item>

        <div class="login-button">
          <el-button type="primary" @click="submit">
            {{ loading ? 'Logining...' : 'Login' }}
          </el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.login-button {
  text-align: center;
  margin-top: 20px;
}
</style>
