<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

import { service } from  "@/js/utils/request.js"
import { validateAccount } from "@/js/utils/validateAccount.js"

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

const addLabel = () => {
  const trimmed = newLabel.value.trim()

  if (trimmed && !form.value.labels.includes(trimmed)) {
    form.value.labels.push(trimmed)
  }

  newLabel.value = '';
}

const submit = async () => {
  try {
    const valid = await formRef.value.validate()
    if (!valid) {
      return;
    }

    console.log(`==> form: ${JSON.stringify(form.value)}`);
    const data = await service.post("/api/v1/auth/account/create_account", form.value);

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
    console.log(`==> error: ${err}`)
  } finally {
    // TODO:
  }
}
</script>

<template>
<el-dialog :model-value="visible" title="Create an account" width="500px">
  <!--hr style="color: #bbb"-->
  <el-divider style="margin: 0 0 20px 0" />

  <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
    <el-form-item label="Firstname" prop="firstname">
      <el-input v-model="form.firstname" placeholder="Enter firstname" />
    </el-form-item>

    <el-form-item label="Lastname" prop="lastname">
      <el-input v-model="form.lastname" placeholder="Enter lastname" />
    </el-form-item>

    <el-form-item label="Email" prop="email">
      <el-input v-model="form.email" placeholder="Enter email" />
    </el-form-item>

    <el-form-item label="Phone" prop="phone">
      <el-input v-model="form.phone" placeholder="Enter phone" />
    </el-form-item>

    <el-form-item label="Password" prop="password">
      <el-input v-model="form.password" type="password" placeholder="Enter password" show-password clearable />
    </el-form-item>

    <!--el-form-item label="Level" prop="level">
      <el-select v-model="form.level" placeholder="Please select a level">
        <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
      </el-select>
    </el-form-item-->

    <el-form :model="form" label-width="100px">
      <div style="display: flex; gap: 5px; align-items: flex-start; flex-wrap: wrap;">
        <el-form-item label="Level" prop="level" style="flex: 1; min-width: 20px;">
          <el-select v-model="form.level" placeholder="Level">
            <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
          </el-select>
        </el-form-item>

        <el-form-item label="Status" prop="status" style="flex: 1; min-width: 20px;">
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
