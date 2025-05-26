<!--ChangePassword :visible="showChangePassword" @close="showChangePassword = false"/-->

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'
import { useRouter } from 'vue-router'

import stores from "@/js/stores"
import { request } from  "@/js/api"

defineProps({
  visible: { type: Boolean },
})

const emit = defineEmits(['update:visible'])

const router = useRouter();


const formRef = ref()

const form = reactive({ oldPassword: '', newPassword: '', confirmPassword: '' });

const validator = (_, value) => {
  if (value !== form.newPassword) {
    return Promise.reject('The passwords do not match twice');
  }

  if (value == form.oldPassword) {
    return Promise.reject('The new password is the same as the old password');
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

const submit = async () => {
  const valid = await formRef.value.validate()
    if (!valid) {
      return
  }

  const loading = ElLoading.service({
    lock: true,
    text: 'Changing password...',
    background: 'rgba(0, 0, 0, 0.7)'
  })

  try {
    const response = await request.post('/api/v1/auth/account/change_password', {
      oldPassword: form.oldPassword,
      newPassword: form.newPassword,
    })

    ElMessage.success('Password changed successfully');
    emit('close');
    stores.clearAccount();
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
<el-dialog
  title="Change password" width="400px"
  :model-value="visible" @update:modelValue="emit('update:visible', $event)"
>
  <el-divider style="margin: 0 0 20px 0" />

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

  <template #footer>
    <el-button @click="emit('update:visible', false)"> Cancel </el-button>
    <el-button type="primary" @click="submit"> Submit </el-button>
  </template>

</el-dialog>
</template>
