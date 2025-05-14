<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

defineProps({
  visible: Boolean,
})

const emit = defineEmits(['close', 'submit'])
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
    { validator: validator, trigger: 'blur' },
  ]
}

const submit = () => {
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
</script>


<template>
<Teleport to="body">
  <div v-if="visible" class="overlay">
    <div class="modal">
      <header class="modal-header"> 🔒 Change Password </header>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="10rem" class="modal-form">
        <el-form-item label="old password" prop="oldPassword">
          <el-input v-model="form.oldPassword" type="password" clearable />
        </el-form-item>

        <el-form-item label="new password" prop="newPassword">
          <el-input v-model="form.newPassword" type="password" clearable />
        </el-form-item>

        <el-form-item label="confirm password" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" clearable />
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
