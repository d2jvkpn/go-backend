<script setup>
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

// import { hello } from "@/js/_hello.js"
// hello()

const allColumns = [
  { prop: 'id',       label: 'ID' },
  { prop: 'username', label: 'Account' },
  { prop: 'email',    label: 'Email' },
  { prop: 'level',    label: 'Level', sortable: true },
  { prop: 'status',   label: 'Status', sortable: true },
]

const visibleColumns = ref(['id', 'username', 'email', 'level', 'status'])
const selectedRows = ref([])

const levels = ['admin', 'editor', 'viewer', 'user', 'guest'];
const statuses = ["activated", "disabled", "deleted"]

const mockAccounts = ref(
  Array.from({ length: 200 }).map((_, i) => ({
    id: i + 1,
    username: `user${i + 1}`,
    email: `user${i + 1}@dev.local`,
    level: levels[i % levels.length],
    status: statuses[i % statuses.length], // i % 2 === 0 ? 'enabled' : 'disabled',
  }))
)

const pagination = ref({ pageIndex: 1, pageSize: 15 })
const filters = ref({ keyword: '', level: '', status: 'activated' })

// 筛选数据
const filteredData = computed(() => {
  return mockAccounts.value.filter(item => {
    const matchKeyword = item.username.includes(filters.value.keyword) || item.email.includes(filters.value.keyword)

    const matchLevel = !filters.value.level || item.level === filters.value.level
    const matchStatus = !filters.value.status || item.status === filters.value.status

    return matchKeyword && matchLevel && matchStatus
  })
})

// 当前页数据
const filteredPagedData = computed(() => {
  const start = (pagination.value.pageIndex - 1) * pagination.value.pageSize
  return filteredData.value.slice(start, start + pagination.value.pageSize)
})

const visibleTableColumns = computed(() =>
  allColumns.filter(col => visibleColumns.value.includes(col.prop))
)

const resetFilters = () => {
  filters.value.keyword = '';
  filters.value.level = '';
  filters.value.status = 'activated';
  pagination.value.pageIndex = 1;
}

const deleteSelected = () => {
  const idsToDelete = selectedRows.value.map(row => row.id)
  // console.log(`~~~ idsToDelete: ${JSON.stringify(idsToDelete)}`)
  let s = idsToDelete.length > 1 ? "s" : ""

  ElMessageBox.confirm(
    `Are you sure you want to delete ${idsToDelete.length} account${s}?`,
    `Delete Account${s} Confirmation`,
    { type: 'warning', confirmButtonText: 'Yes', cancelButtonText: 'No' }
  )
  .then(() => {
    mockAccounts.value = mockAccounts.value.filter(user => !idsToDelete.includes(user.id))
    selectedRows.value = []
    ElMessage.success(`Deleted ${idsToDelete.length} account${s}`)
  })
  .catch(() => {
    ElMessage.info(`Delete account${s} canceled.`)
  })
}

const handleSizeChange = (size) => {
  pagination.value.pageIndex = 1
  pagination.value.pageSize = size
}

function customSort (a, b) {
  return a > b;
}
</script>


<template>
<div class="toolbar"> <!-- 搜索栏 -->
  <div class="toolbar-left">  <!-- 左侧：搜索、角色、重置 -->
    <el-input v-model="filters.keyword" placeholder="seach account or email" style="width: 200px" clearable />

    <el-select v-model="filters.level" placeholder="level" style="width: 150px" clearable>
      <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
    </el-select>

    <el-select v-model="filters.status" placeholder="status" style="width: 150px" clearable>
      <el-option v-for="e in statuses" :value="e" :label="e" :key="`account::status::${e}`" />
    </el-select>

    <el-button @click="resetFilters"> Reset </el-button>
  </div>

  <div class="toolbar-right">
    <el-dropdown trigger="click">
      <el-button type="primary">
        Columns <el-icon> <ArrowDown /> </el-icon>
      </el-button>

      <template #dropdown>
        <el-dropdown-menu class="column-dropdown">
          <el-checkbox-group v-model="visibleColumns">
            <el-dropdown-item v-for="col in allColumns" :key="col.prop" class="no-hover">
              <el-checkbox :value="col.prop"> {{ col.label }} </el-checkbox>
            </el-dropdown-item>
          </el-checkbox-group>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <el-button type="danger" @click="deleteSelected" :disabled="!selectedRows.length">
      Delete
    </el-button>
  </div>
</div>

<el-table :data="filteredPagedData" style="margin-top: 10px;"  @selection-change="selectedRows = $event" border>
  <el-table-column type="selection" width="50" />

  <!--el-table-column prop="id" label="ID" sortable :sort-method="customSort"/-->

  <el-table-column v-for="col in visibleTableColumns"
    :prop="col.prop" :label="col.label" :key="col.prop" :sortable="col.sortable"
  />
</el-table>


<el-pagination
  v-model:current-page="pagination.pageIndex"
  v-model:page-size="pagination.pageSize"
  :total="filteredData.length"
  :page-sizes="[15, 20, 50]"
  class="pagination"
  @size-change="handleSizeChange"
  layout="prev, pager, next, jumper, total, sizes"
/>
</template>


<style scoped>
.toolbar {
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.toolbar-left, .toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.column-dropdown {
  padding: 8px 16px;
}

.column-dropdown .no-hover {
  background-color: transparent !important;
  cursor: default;
  padding: 0;
}

.pagination {
  margin-top: 16px;
  text-align: right;
  display: flex;
  justify-content: flex-end;
}
</style>
