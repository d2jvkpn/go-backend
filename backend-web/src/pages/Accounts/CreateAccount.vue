<script setup>
import { ref } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'

import { request } from "@/js/api";
import { validateAccount, validateContact } from "@/js/utils/validateAccount.js"

defineProps({
  visible: { type: Boolean },
})

const emit = defineEmits(['update:visible', 'success'])

const formRef = ref()

const form = ref({
  firstname: '',
  lastname: '',
  email: '',
  phone: '',
  password: '',
  level: '',
  labels: [],
  status: 'activated',
})

const levels = ['admin', 'editor', 'reviewer', 'user', 'guest'];

const newLabel = ref('')

function addLabel () {
  const trimmed = newLabel.value.trim()

  if (trimmed && !form.value.labels.includes(trimmed)) {
    form.value.labels.push(trimmed)
  }

  newLabel.value = '';
}

const submit = async () => {
  if (!validateContact(form.value)) {
    return
  }

  const valid = await formRef.value.validate()
  if (!valid) {
    return;
  }

  const loading = ElLoading.service({
    lock: true,
    text: 'Creating account...',
    background: 'rgba(0, 0, 0, 0.7)',
  })
  console.log(`==> form: ${JSON.stringify(form.value)}`);

  try {
    const data = await request.post("/api/v1/auth/account/create_account", form.value);

    form.value = {
      firstname: '',
      lastname: '',
      email: '',
      phone: '',
      password: '',
      level: '',
      labels: [],
      status: 'activated',
    }

    ElMessage.success(`The account has been successfully created: ${data.accountId}`);
    emit('success');
    emit('update:visible', false);
  } catch (err) {
    console.log(`==> CreateAccount error: ${err}`)
  } finally {
    loading?.close()
  }
}
</script>

<template>
<el-dialog :model-value="visible" title="Create an account" width="500px">
  <!--hr style="color: #bbb"-->
  <el-divider style="margin: 0 0 20px 0" />

  <el-form :model="form" :rules="validateAccount" ref="formRef" label-width="100px">
    <el-form-item label="Firstname" prop="firstname" required>
      <el-input v-model="form.firstname" placeholder="Enter firstname" />
    </el-form-item>

    <el-form-item label="Lastname" prop="lastname" required>
      <el-input v-model="form.lastname" placeholder="Enter lastname" />
    </el-form-item>

    <el-form-item label="Email" prop="email">
      <el-input v-model="form.email" placeholder="Enter email" />
    </el-form-item>

    <el-form-item label="Phone" prop="phone">
      <el-input v-model="form.phone" placeholder="Enter phone" />
    </el-form-item>

    <el-form-item label="Password" prop="password" required>
      <el-input v-model="form.password" type="password" placeholder="Enter password" show-password clearable />
    </el-form-item>

    <!--el-form-item label="Level" prop="level">
      <el-select v-model="form.level" placeholder="Please select a level">
        <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
      </el-select>
    </el-form-item-->

    <el-form :model="form" label-width="100px">
      <div style="display: flex; gap: 5px; align-items: flex-start; flex-wrap: wrap;">
        <el-form-item label="Level" prop="level" style="flex: 1; min-width: 20px;" required>
          <el-select v-model="form.level" placeholder="Level">
            <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
          </el-select>
        </el-form-item>

        <el-form-item label="Status" prop="status" style="flex: 1; min-width: 20px;" required>
          <el-select v-model="form.status" placeholder="Status">
            <el-option value="activated" label="activated" />
            <el-option value="created" label="created" />
          </el-select>
        </el-form-item>
      </div>
    </el-form>

    <el-form-item label="Labels" prop="labels">
      <div style="display: flex; flex-wrap: wrap; gap: 8px;">
        <el-tag v-for="(label, index) in form.labels" :key="label" @close="form.labels.splice(index, 1)" closable>
          {{ label }}
        </el-tag>
        <el-input
          v-model="newLabel" @keyup.enter="addLabel"
          placeholder="input label and press enter" size="small" style="width: 180px; height: 30px"
        />
      </div>
    </el-form-item>
  </el-form>

  <template #footer>
    <el-button @click="$emit('update:visible', false)"> Cancel </el-button>
    <el-button type="primary" @click="submit"> Submit </el-button>
  </template>
</el-dialog>
</template>
