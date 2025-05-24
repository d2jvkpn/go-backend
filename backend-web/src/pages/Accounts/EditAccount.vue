<script setup>
import { computed, reactive, watch, ref } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'
import cloneDeep from 'lodash/cloneDeep'

import { service } from "@/js/utils/request.js"
import { validateAccount } from "@/js/utils/validateAccount.js"

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
