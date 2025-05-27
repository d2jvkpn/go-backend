<script setup>
import { ref, onBeforeMount, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { request } from  "@/js/api"
import { getFirstRoute } from "@/router/index"
import stores from "@/js/stores"

// console.log(`==> import.meta.env: ${JSON.stringify(import.meta.env)}`);

const account = ref('')
const password = ref('')
const captcha = ref({ enabled: true, id: "", base64Image: "", length: 0, answer: "" })
const loading = ref(false)
const router = useRouter()

let firstRoute = ""

async function getCaptcha() {
  try {
    const data = await request.get("/api/v1/open/account/captcha");
    captcha.value.enabled = data.enabled;
    captcha.value.id = data.id;
    captcha.value.base64Image = data.base64Image;
    captcha.value.length = data.length;
  } catch (err) {
    console.log(`!!! getCaptcha: ${err}`)
  }
}

async function submit () {
  if (loading.value) {
    ElMessage.warning('Logining...');
    return;
  }

  if (!account.value || !password.value) {
    // alert('Please enter acocunt and password!')
    ElMessage.warning('Please enter acocunt and password!');
    return
  }

  if (captcha.value.enabled) {
    if (captcha.value.answer.length == 0) {
      ElMessage.warning('Please enter captcha!');
      return
    }

    console.log("???", captcha.value.answer.length, captcha.value.length)
    if (captcha.value.answer.length != captcha.value.length) {
      ElMessage.warning('Invlaid captcha!');
      return
    }
  }

  const loginData = { password: password.value }
  if (account.value.includes("@")) {
    loginData.email = account.value;
  } else {
    loginData.phone = account.value;
  }

  if (captcha.value.enabled) {
    loginData.captchaId = captcha.value.id
    loginData.captchaAnswer = captcha.value.answer
  }

  try {
    const data = await request.post(
      `${import.meta.env.VITE_API_URL}/api/v1/open/account/login`,
      loginData,
      { params: {platform: "web"} },
    )

    firstRoute = getFirstRoute(data.level);

    stores.setAccount(data)
    ElMessage.success(`Welcome back, ${data.firstname} ${data.lastname}!`)

    router.push(firstRoute)
  } catch (err) {
    console.log(`!!! Login error: ${err.message}`)
    if (err.code == "captcha_verify_failed") {
      await getCaptcha();
    }
  } finally {
    loading.value = false
  }
}

onBeforeMount(() => {
  if (!stores.checkIsLoggedIn()) {
    return
  }

  const account = stores.getAccount();
  if (account?.level) {
    router.push(getFirstRoute(account.level));
  }
})

onMounted(async () => {
  await getCaptcha();
})

// onBeforeMount, onMounted, onBeforeUpdate, onUpdated, onUnmounted
</script>


<template>
<div class="login">
  <el-card class="login-card">
    <template #header> Please Login </template>
    <el-form @submit.prevent="submit">
      <el-form-item>
        <el-input placeholder="Enter Email or Phone" v-model="account" />
      </el-form-item>

      <el-form-item>
        <el-input placeholder="Enter Password" type="password" v-model="password" clearable show-password/>
      </el-form-item>

      <el-form-item v-if="captcha.enabled">
        <div class="captcha-row">
          <el-input placeholder="Enter CAPTCHA" maxlength="6" class="captcha-answer" v-model="captcha.answer"/>
          <img class="captcha-img" alt="Captcha Image" :src="captcha.base64Image" @click="getCaptcha" />
        </div>
      </el-form-item>

      <div class="login-button">
        <el-button type="primary" @click="submit"> {{ loading ? 'Logining...' : 'Login' }} </el-button>
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
  width: 24rem;
  margin: 0rem 2rem;
}

.login-button {
  display: flex;
  justify-content: center;
  align-items: center;
}

.captcha-row {
  display: flex;
  width: 100%;
  gap: 12px;
  align-items: center;
}

.captcha-answer {
  flex: 0.5;
}

.captcha-img {
  height: 32px;
  flex: 0.5;
  cursor: pointer;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}
</style>
