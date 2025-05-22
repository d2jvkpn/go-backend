<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'
import { useRouter } from 'vue-router'

import { clearAccount } from "@/js/stores/local.js"
import { service } from  "@/js/utils/request.js"

defineProps({
  visible: { type: Boolean },
})

const router = useRouter();

const emit = defineEmits(['close'])
const formRef = ref()

const form = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' })

const validator = (_, value) => {
  if (value !== form.newPassword) {
    return Promise.reject('The passwords do not match twice')
  }

  if (value == form.oldPassword) {
    return Promise.reject('The new password is the same as the old password')
  }

  return Promise.resolve()
}

const rules = {
  oldPassword: [{ required: true, message: 'Please enter old password', trigger: 'blur' }],
  newPassword: [{ required: true, message: 'Please enter new password', trigger: 'blur' }],
  confirmPassword: [
    { required: true, message: 'Please confirm password', trigger: 'blur' },
    {
      validator: validator,
      trigger: 'blur',
      min: 8,
      max: 32,
      pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[^]{8,32}$/,
      message: 'Password must be 8-32 chars with at least one uppercase, lowercase and number'
    },
  ]
}

const submitV1 = () => {
  formRef.value.validate((valid) => {
    if (!valid) {
      return
    }

    // emit('submit', { ...form })
    console.log(`--> change password: ${JSON.stringify(form)}`)
    ElMessage.success('The password has been successfully changed')
    emit('close')
  })
}

const submit = async () => {
  const loading = ElLoading.service({
    lock: true,
    text: 'Changing password...',
    background: 'rgba(0, 0, 0, 0.7)'
  })

  try {
    const valid = await formRef.value.validate()
    if (!valid) {
      return
    }
    emit('close');

    const response = await service.post('/api/v1/auth/account/change_password', {
      oldPassword: form.oldPassword,
      newPassword: form.newPassword
    })

    ElMessage.success('Password changed successfully')
    clearAccount();
    router.push('/login');
  } catch (err) {
    console.log(`!!! error: ${JSON.stringify(err)}, ${err.message}`)
    ElMessage.error(err.message);
  } finally {
    loading?.close()
  }
}
</script>


<template>
<Teleport to="body">
  <div v-if="visible" class="overlay">
    <div class="modal">
      <header class="modal-header"> 🔒 Change Password </header>
      <hr style="color: #bbb">

      <el-form ref="formRef" :model="form" :rules="rules" label-width="10rem" class="modal-form">
        <el-form-item label="old password" prop="oldPassword">
          <el-input v-model="form.oldPassword" type="password" show-password clearable/>
        </el-form-item>

        <el-form-item label="new password" prop="newPassword">
          <el-input v-model="form.newPassword" type="password" show-password clearable/>
        </el-form-item>

        <el-form-item label="confirm password" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password clearable/>
        </el-form-item>
      </el-form>

      <footer class="modal-footer">
        <el-button @click="$emit('close')"> Cancel </el-button>
        <el-button type="primary" @click="submit"> Submit </el-button>
      </footer>
    </div>
  </div>
</Teleport>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  background: white;
  border-radius: 5px;
  width: 30rem;
  padding: 2rem;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

.modal-header {
  font-weight: bold;
  font-size: 1.2rem;
  margin-bottom: 1rem;
}

.modal-form {
  margin-top: 1rem;
}

.modal-footer {
  margin-top: 2rem;
  display: flex;
  justify-content: flex-end;
  gap: 2rem;
}
</style>
