<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

// import { hello } from "@/js/_hello.js"
// hello()
import { service } from "@/js/utils/request.js";
import CreateAccount from './CreateAccount.vue'
import UpdateStatus from "./UpdateStatus.vue";


// create an account
const createAccountVisible = ref(false)

function openCreateAccount() {
  console.log("==> openCreateAccount")
  createAccountVisible.value = true
}

// show accounts table
const allColumns = [
  { prop: 'id',        label: 'ID' },
  { prop: 'firstname', label: 'Firstname' },
  { prop: 'lastname',  label: 'Lastname' },
  { prop: 'email',     label: 'Email' },
  { prop: 'phone',     label: 'Phone', width: 120 },
  { prop: 'level',     label: 'Level', width: 100 },
  { prop: 'labels',    label: 'Labels', width: 250 },
  { prop: 'createdAt', label: 'Created At', sortable: true },
  { prop: 'updatedAt', label: 'updated At', sortable: true },
  //{ prop: 'status',    label: 'Status' },
]

// const visibleColumns = ref(['id', 'firstname', 'lastname', 'email', 'phone', 'level', 'labels', 'createdAt', 'status'])
const visibleColumns = ref(['email', 'phone', 'level', 'labels', 'createdAt'])

const selectVisibleColumns = computed(() =>
  allColumns.filter(col => visibleColumns.value.includes(col.prop))
)

// fetch accounts
const pageData = ref({ total: 0, items: [] })
const loading = ref(false)
const error = ref(null);

const fetchData = async () => {
  loading.value = true;
  error.value = null;

  try {
     let data = await service.get("/api/v1/auth/account/query_accounts", {}, { params: query.value });

     if (query.value.pageIndex == 1) {
       pageData.value.total = data.total;
     }

     data.items.forEach(item => {
       item.createdAt = dayjs(item.createdAt).format('YYYY-MM-DD HH:mm')
       item.updatedAt = dayjs(item.updatedAt).format('YYYY-MM-DD HH:mm')
    })

     pageData.value.items = data.items
  } catch (err) {
    error.value = err.response?.data?.msg || err.msg;
  } finally {
    loading.value = false;
  }
};

// search bar
const levels = ['admin', 'editor', 'reviewer', 'user', 'guest'];
const statuses = ["created", "activated", "blocked"]

const query = ref({ pageIndex: 1, pageSize: 10, keyword: '', level: '', status: '' })

const handleReset = async () => {
  query.value.pageIndex = 1;
  query.value.keyword = '';
  query.value.level = '';
  query.value.status = '';
  await fetchData();
}

const handleSearch = async () => {
  query.value.pageIndex = 1;
  await fetchData();
};

const updatePageIndex = async (v) => {
  console.log(`==> updatePageIndex: ${v}`)
  query.value.pageIndex = v;
  await fetchData();;
}

const updatePageSize = async (v) => {
  console.log(`==> updatePageSize: ${v}`)
  query.value.pageIndex = 1;
  query.value.pageSize = v;
  await fetchData();
}

function customSort (a, b) {
  return a > b;
}

//
const selectedRows = ref([]);

const deleteSelected = async () => {
  const idsToDelete = selectedRows.value.map(row => row.id)
  // console.log(`~~~ idsToDelete: ${JSON.stringify(idsToDelete)}`)
  let s = idsToDelete.length > 1 ? "s" : ""

  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete ${idsToDelete.length} account${s}?`,
      `Delete Account${s} Confirmation`,
      { type: 'warning', confirmButtonText: 'Yes', cancelButtonText: 'No' }
    );

    const data = await service.post(
      "/api/v1/auth/account/delete_accounts",
      {},
      { params: { "accountId": idsToDelete } },
    );

    pageData.value.items = pageData.value.items.filter(item => !idsToDelete.includes(item.id));
    ElMessage.success(`Deleted ${data.count}/${idsToDelete.length} account${s}`);
  } catch (err) {
    ElMessage.info(`Delete account${s} canceled.`)
  }
}

//
const editAccount = (data) => {
  console.log(`==> ${data.id}: ${data.firstname} ${data.lastname}, ${data.status}`)
}

//
const updateStatusVisible = ref(false)
const selectedAccount = ref(null)

function openUpdateStatus(account) {
  selectedAccount.value = account
  updateStatusVisible.value = true
}

function updateStatus (account) {
  console.log(`==> ${account.id}: ${account.firstname} ${account.lastname}, ${account.status}`)
  selectedAccount.value = account
  updateStatusVisible.value = true
}

//
watch(
  // query.value.keyword
  () => [ query.value.pageSize, query.value.pageIndex, query.value.level, query.value.status ],
  (newValues, oldValues) => {
    console.log(`==> watch: ${newValues}, ${oldValues}`);
    fetchData();
  }
);

onMounted(async () => {
  console.log("==> onMounted");
  await fetchData();
});
</script>


<template>
<div class="toolbar"> <!-- 搜索栏 -->
  <div class="toolbar-left">  <!-- 左侧：搜索、角色、重置 -->
    <el-input
      v-model="query.keyword"
      placeholder="seach account or email"
      @keyup.enter="handleSearch"
      style="width: 200px"
      clearable
    />

    <el-select v-model="query.level" placeholder="level" style="width: 100px" clearable>
      <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
    </el-select>

    <el-select v-model="query.status" placeholder="status" style="width: 100px" clearable>
      <el-option v-for="e in statuses" :value="e" :label="e" :key="`account::status::${e}`" clearable/>
    </el-select>

    <el-button type="info" @click="handleReset"> Reset </el-button>
    <!--el-button type="info" @click="handleSearch"> Search </el-button-->
  </div>

  <div class="toolbar-right">
    <el-button type="danger" @click="deleteSelected" :disabled="!selectedRows.length">
      Delete
    </el-button>

    <el-dropdown trigger="click">
      <el-button type="success">
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

    <el-button type="primary" @click="openCreateAccount">Create</el-button>
    <!--CreateAccount v-model:visible="openCreateAccount" @success="refreshData" /-->
    <CreateAccount v-model:visible="createAccountVisible" @success="handleSearch" />
  </div>
</div>

<UpdateStatus
  v-model:visible="updateStatusVisible"
  :account="selectedAccount"
  @update:status="({ id, status }) => {
    const target = pageData.items.find(item => item.id === id)
    if (target) { target.status = status }
  }"
/>

<el-table
  :data="pageData.items"
  style="margin-top: 10px;"
  @selection-change="selectedRows = $event"
  border
  v-loading="loading"
  empty-text="No accounts found"
>
  <el-table-column type="selection" width="40" />

  <el-table-column label="Full Name" prop="fullName" width="120">
    <template #default="scope"> {{ scope.row.firstname }} {{ scope.row.lastname }} </template>
  </el-table-column>

  <!--el-table-column prop="id" label="ID" sortable :sort-method="customSort"/-->
  <el-table-column v-for="col in selectVisibleColumns"
    :prop="col.prop" :label="col.label" :key="col.prop" :sortable="col.sortable" :width="col.width"
  />

  <el-table-column label="Actions" fixed="right" width="180">
    <template #default="scope">
      <el-button type="warning" size="small" @click="updateStatus(scope.row)"> {{ scope.row.status }} </el-button>
      <el-button type="primary" size="small" @click="editAccount(scope.row)"> edit </el-button>
    </template>
  </el-table-column>
</el-table>

<el-alert v-if="error" :title="error" type="error" show-icon style="margin-top: 10px" />

<el-pagination
  v-model:current-page="query.pageIndex"
  v-model:page-size="query.pageSize"
  :total="pageData.total"
  :page-sizes="[10, 20, 50]"
  class="pagination"
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
