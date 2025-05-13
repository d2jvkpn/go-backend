<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

// import { hello } from "@/js/utils/hello.js"
// hello()

const allColumns = [
  { prop: 'id', label: 'ID' },
  { prop: 'username', label: 'Account' },
  { prop: 'email', label: 'Email' },
  { prop: 'role', label: 'Role' },
  { prop: 'status', label: 'Status' },
]

const visibleColumns = ref(['id', 'username', 'email', 'role'])
const selectedRows = ref([])

const mockAccounts = ref(
  Array.from({ length: 100 }).map((_, i) => ({
    id: i + 1,
    username: `user${i + 1}`,
    email: `user${i + 1}@example.com`,
    role: ['admin', 'normal', 'visitor'][i % 3],
    status: i % 2 === 0 ? 'enabled' : 'disabled',
  }))
)

const pagination = ref({ currentPage: 1, pageSize: 15 })
const filters = ref({ keyword: '', role: '' })

// 筛选数据
const filteredData = computed(() => {
  return mockAccounts.value.filter(item => {
    const matchKeyword = item.username.includes(filters.value.keyword) || item.email.includes(filters.value.keyword)

    const matchRole = !filters.value.role || item.role === filters.value.role

    return matchKeyword && matchRole
  })
})

// 当前页数据
const filteredPagedData = computed(() => {
  const start = (pagination.value.currentPage - 1) * pagination.value.pageSize
  return filteredData.value.slice(start, start + pagination.value.pageSize)
})

const visibleTableColumns = computed(() =>
  allColumns.filter(col => visibleColumns.value.includes(col.prop))
)

const resetFilters = () => {
  filters.value.keyword = '';
  filters.value.role = '';
  pagination.value.currentPage = 1;
}

const deleteSelected = () => {
  const idsToDelete = selectedRows.value.map(row => row.id)
  mockAccounts.value = mockAccounts.value.filter(user => !idsToDelete.includes(user.id))
  selectedRows.value = []
  ElMessage.success(`Deleted ${idsToDelete.length} accounts`)
}

const handleSizeChange = (size) => {
  pagination.value.currentPage = 1
  pagination.value.pageSize = size
}
</script>


<template>
<div class="toolbar"> <!-- 搜索栏 -->
  <div class="toolbar-left">  <!-- 左侧：搜索、角色、重置 -->
    <el-input
      v-model="filters.keyword"
      placeholder="seach account or email"
      clearable
      style="width: 200px"
    />

    <el-select v-model="filters.role" placeholder="role" clearable style="width: 150px">
      <el-option label="admin" value="admin" />
      <el-option label="normal" value="normal" />
      <el-option label="visitor" value="visitor" />
    </el-select>
  <el-button @click="resetFilters">Reset</el-button>
  </div>

  <div class="toolbar-right">
    <el-dropdown trigger="click">
      <el-button type="primary">
        Columns
        <el-icon> <ArrowDown /> </el-icon>
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

<el-table
  :data="filteredPagedData"
  border
  style="margin-top: 10px;"
  @selection-change="selectedRows = $event"
>
  <el-table-column type="selection" width="50" />
  <el-table-column
    v-for="col in visibleTableColumns"
    :key="col.prop"
    :label="col.label"
    :prop="col.prop"
  />
</el-table>


<el-pagination
  v-model:current-page="pagination.currentPage"
  v-model:page-size="pagination.pageSize"
  :total="filteredData.length"
  :page-sizes="[15, 20, 30]"
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
