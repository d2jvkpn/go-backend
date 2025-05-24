<script setup>
import { computed, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus';

import { service } from "@/js/utils/request.js"

const props = defineProps({
  visible: Boolean,
  account: Object,
})

const emit = defineEmits(['update:visible', 'update:status'])

const form = reactive({
  id: '',
  firstname: '',
  lastname: '',
  level: '',
  status: '',
  newStatus: '',
})

const statusOptionsMap = {
  created: ['activated', 'deleted'],
  activated: ['blocked', 'deleted'],
  blocked: ['activated', 'deleted'],
  deleted: []
}

const statusOptions = computed(() => {
  return statusOptionsMap[form.status] || []
})

async function confirm() {
  if (!form.newStatus || form.newStatus === form.status) {
    ElMessage.warning('Please select a different status')
    return
  }

  console.log(`==> update status: id=${form.id}, status=${form.status}, newStatus=${form.newStatus}`)

  try {
    await service.post("/api/v1/auth/account/update_status", {}, { params: {
      accountId: form.id,
      status: form.status,
      newStatus: form.newStatus,
    } })

    emit('update:status', { id: form.id, status: form.newStatus })
    emit('update:visible', false)
    ElMessage.success('Status updated successfully')
  } catch (err) {
    //
  } finally {
    //
  }


}

watch(() => props.account, (newVal) => {
  if (newVal) {
    Object.assign(form, {
      id: newVal.id,
      firstname: newVal.firstname,
      lastname: newVal.lastname,
      level: newVal.level,
      status: newVal.status,
      newStatus: ''
    })
  }
}, { immediate: true })
</script>

<template>
<el-dialog
  :model-value="visible" title="Update Account" width="400px"
  @update:modelValue="emit('update:visible', $event)"
>
  <el-divider style="margin: 0 0 20px 0" />

  <el-form :model="form" label-width="120px">
    <el-form-item label="Account ID"> <el-input :value="form.id" disabled /> </el-form-item>

    <el-form-item label="Full Name">
      <el-input :value="form.firstname + ' ' + form.lastname" disabled />
    </el-form-item>

    <el-form-item label="Level"> <el-input :value="form.level" disabled /> </el-form-item>

    <el-form-item label="Current Status">
      <el-tag :type="form.status === 'activated' ? 'success' : 'info'" size="large">
        {{ form.status }}
      </el-tag>
    </el-form-item>

    <el-form-item label="Update Status" v-if="statusOptions.length > 0">
      <div style="display: flex; gap: 12px;">
        <el-checkbox
          v-for="val in statusOptions"
          :key="val" :label="val"
          :model-value="form.newStatus === val"
          @change="() => form.newStatus = val"
        />
      </div>
    </el-form-item>
  </el-form>

  <template #footer>
    <el-button @click="emit('update:visible', false)"> Cancel </el-button>
    <el-button type="primary" @click="confirm"> Confirm </el-button>
  </template>
</el-dialog>
</template>
