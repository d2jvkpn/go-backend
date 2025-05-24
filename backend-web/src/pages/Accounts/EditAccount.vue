<script setup>
import { computed, reactive, watch, ref } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'
import cloneDeep from 'lodash/cloneDeep'

import { service } from "@/js/utils/request.js"

const props = defineProps({
  visible: Boolean,
  account: Object,
})

const emit = defineEmits(['update:visible', 'update:refresh'])

//
const formRef = ref()

const form = reactive({
  id: '',
  firstname: '',
  lastname: '',
  email: '',
  phone: '',
  level: '',
  labels: [],

  password: '',
})

const levels = ['admin', 'editor', 'reviewer', 'user', 'guest'];

const newLabel = ref('')

function addLabel ()  {
  const trimmed = newLabel.value.trim()

  if (trimmed && !form.labels.includes(trimmed)) {
    form.labels.push(trimmed)
  }

  newLabel.value = '';
}

const validateContact = (rule, value, callback) => {
  const email = form.value.email?.trim()
  const phone = form.value.phone?.trim()
  if (!email && !phone) {
    // Please enter at least an email or a phone number
    callback(new Error('Either email or phone must be provided'))
  } else {
    callback()
  }
}

const validateLabels = (rule, value, callback) => {
  if (!Array.isArray(value)) {
    callback();
    return;
  }

  const tooLong = value.find(label => label.length > 32);
  if (tooLong) {
    callback(new Error(`Label ${tooLong} exceeds 32 characters`));
    return;
  }

  if (value.length > 16) {
    callback(new Error('You can select up to 16 labels only'));
    return;
  }

  callback();
}

const rules = {
  firstname: [{ required: true, min: 2, max: 32, message: 'Please enter firstname', trigger: 'blur' }],
  lastname: [{ required: true, min: 2, max: 32, message: 'Please enter lastname', trigger: 'blur' }],
  email: [
    { validator: validateContact, trigger: 'blur' },
    { type: 'email', message: 'invalid email format', trigger: 'blur' },
    { min: 5, max: 64, message: 'Email too long', trigger: 'blur' },
  ],
  phone: [{ min: 6, max: 20, validator: validateContact, trigger: 'blur' }],
  password: [{
    required: true, message: 'Please enter password',
    trigger: 'blur', min: 8, max: 32,
    pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[^]{8,32}$/,
    message: 'Password must be 8-32 chars with at least one uppercase, lowercase and number',
  }],
  level: [{ required: true, message: 'Please enter level', trigger: 'change' }],
  labels: [{ validator: validateLabels, trigger: 'change' }],
}

async function confirm() {
  const valid = await formRef.value.validate()
  if (!valid) {
    return;
   }

  console.log(`==> update account`)

  const loading = ElLoading.service({
    lock: true,
    text: 'Updating account...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    await service.post("/api/v1/auth/account/edit_account", form, { params: { accountId: form.id } })

    emit('update:visible', false)
    emit('update:refresh', form)
    ElMessage.success('Account updated successfully')
  } catch (err) {
    console.log(`!!! edit status error: ${err}`)
  } finally {
    loading?.close()
  }
}

watch(() => props.visible, (visible) => {
  if (visible && props.account) {
    const cleaned = cloneDeep(props.account)
    delete cleaned.createdAt
    delete cleaned.updatedAt
    delete cleaned.status

    cleaned.labels = Array.isArray(cleaned.labels) ? cleaned.labels : []
    Object.assign(form, cleaned)

    newLabel.value = ''
  }
}, { immediate: true })

</script>

<template>
<el-dialog
  :model-value="visible" title="Update account" width="500px"
  @update:modelValue="emit('update:visible', $event)"
>
  <el-divider style="margin: 0 0 20px 0" />

  <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
    <el-form-item label="ID"> <el-input :value="form.id" disabled /> </el-form-item>

    <el-form-item label="Firstname"> <el-input v-model="form.firstname"/> </el-form-item>
    <el-form-item label="Lastname"> <el-input v-model="form.lastname"/> </el-form-item>
    <el-form-item label="Email"> <el-input v-model="form.email"/> </el-form-item>
    <el-form-item label="Phone"> <el-input v-model="form.phone"/> </el-form-item>

    <el-form-item label="Level" prop="level" style="flex: 1; min-width: 20px;">
      <el-select v-model="form.level" placeholder="Level">
        <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
      </el-select>
    </el-form-item>

    <el-form-item label="Labels" prop="labels">
      <div style="display: flex; flex-wrap: wrap; gap: 8px;">
        <el-tag
          v-for="(label, index) in form.labels"
          :key="label"
          @close="form.labels.splice(index, 1)"
          closable
        >
          {{ label }}
        </el-tag>
        <el-input
          v-model="newLabel" @keyup.enter="addLabel"
          placeholder="input label and press enter" size="small" style="width: 180px; height: 30px"
        />
      </div>
    </el-form-item>

    <el-form-item label="Password">
      <el-input v-model="form.password" type="password" placeholder="Enter password" show-password clearable />
    </el-form-item>
  </el-form>

  <template #footer>
    <el-button @click="emit('update:visible', false)"> Cancel </el-button>
    <el-button type="primary" @click="confirm"> Confirm </el-button>
  </template>
</el-dialog>
</template>
