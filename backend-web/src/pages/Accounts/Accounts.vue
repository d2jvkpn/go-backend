<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
// https://element-plus.org/zh-CN/component/icon.html
import { ArrowDown, Edit, CopyDocument } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import * as yaml from 'js-yaml'

import { request } from "@/js/api";
import CreateAccount from './CreateAccount.vue'
import UpdateStatus from "./UpdateStatus.vue";
import EditAccount from "./EditAccount.vue";


// create an account
const createAccountVisible = ref(false)

function openCreateAccount() {
  console.log("==> openCreateAccount")
  createAccountVisible.value = true
}

function formatTime(row, column, cellValue, index) {
  return dayjs(cellValue).format('YYYY-MM-DD HH:mm')
}

// show accounts table
const allColumns = [
  { prop: 'id',        label: 'ID' },
  { prop: 'firstname', label: 'Firstname' },
  { prop: 'lastname',  label: 'Lastname' },
  { prop: 'email',     label: 'Email' },
  { prop: 'phone',     label: 'Phone', width: 120 },
  { prop: 'level',     label: 'Level', width: 100 },
  { prop: 'labels',    label: 'Labels' },
  { prop: 'createdAt', label: 'Created At', sortable: true, width: 150, formatter: formatTime },
  { prop: 'updatedAt', label: 'Updated At', sortable: true, width: 150, formatter: formatTime },
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

async function fetchData() {
  loading.value = true;
  error.value = null;

  try {
     const [sortBy, order] = sortValue.value.split('-');
     const params = {...query.value, sortBy, order};
     const data = await request.get("/api/v1/auth/account/query_accounts", { params });

     if (query.value.pageIndex == 1) {
       pageData.value.total = data.total;
     }

     /*
     data.items.forEach(item => {
       item._createdAt = dayjs(item.createdAt).format('YYYY-MM-DD HH:mm')
       item._updatedAt = dayjs(item.updatedAt).format('YYYY-MM-DD HH:mm')
    })
    */

     pageData.value.items = data.items
  } catch (err) {
    error.value = err.response?.data?.msg || err.msg;
  } finally {
    loading.value = false;
  }
};

// search bar
const levels = ['admin', 'editor', 'reviewer', 'user', 'guest'];
const statuses = ["created", "activated", "blocked", "deleted"]

const query = ref({ pageIndex: 1, pageSize: 10, keyword: '', level: '', status: '' })

const sortValue = ref('createdAt-desc')

async function handleReset () {
  query.value.pageIndex = 1;
  query.value.keyword = '';
  query.value.level = '';
  query.value.status = '';
  sortValue.value = 'createdAt-desc';

  await fetchData();
}

async function handleSearch () {
  query.value.pageIndex = 1;
  await fetchData();
};

async function updatePageIndex (v) {
  console.log(`==> updatePageIndex: ${v}`)
  query.value.pageIndex = v;
  await fetchData();
}

async function updatePageSize (v) {
  console.log(`==> updatePageSize: ${v}`)
  query.value.pageIndex = 1;
  query.value.pageSize = v;
  await fetchData();
}

//
const selectedRows = ref([]);

async function deleteSelected () {
  const idsToDelete = selectedRows.value.map(row => row.id)
  // console.log(`~~~ idsToDelete: ${JSON.stringify(idsToDelete)}`)
  let s = idsToDelete.length > 1 ? "s" : ""

  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete ${idsToDelete.length} account${s}?`,
      `Delete Account${s} Confirmation`,
      { type: 'warning', confirmButtonText: 'Yes', cancelButtonText: 'No' }
    );

    const data = await request.post(
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
const selectedAccount = ref(null)
const updateStatusVisible = ref(false)

function updateStatus (account) {
  console.log(`==> UpdateStatus: ${account.id}: ${account.firstname} ${account.lastname}, ${account.status}`)
  selectedAccount.value = account
  updateStatusVisible.value = true
}

//
const editAccountVisible = ref(false)

function editAccount (account) {
  console.log(`==> EditAccount: ${account.id}, ${account.firstname} ${account.lastname}, ${account.status}`)
  selectedAccount.value = account
  editAccountVisible.value = true
}

//
function copyToClipboard(obj) {
  // const text = JSON.stringify(obj, null, 2)
  const text = yaml.dump(obj);

  navigator.clipboard.writeText(text)
  .then(() => {
    ElMessage.success(`Copy success: ${obj.firstname} ${obj.lastname}`);
  })
  .catch(err => {
    ElMessage.error(`Copy failed: ${err}`);
  })
}

//
watch(
  // query.value.keyword
  () => [query.value.pageSize, query.value.level, query.value.status, sortValue.value ],
  (newValues, oldValues) => {
    console.log(`==> Watch: ${newValues}, ${oldValues}`);
    query.value.pageIndex = 1;
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
      placeholder="Seaching" style="width: 240px"
      title="*Firstname*, *Lastname*, *Email*, *Phone* and ^Labels$"
      v-model="query.keyword"  @keyup.enter="handleSearch" clearable
    />

    <el-select placeholder="level" style="width: 100px" v-model="query.level" clearable>
      <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
    </el-select>

    <el-select style="width: 100px" placeholder="status" v-model="query.status" clearable>
      <el-option v-for="e in statuses" :value="e" :label="e" :key="`account::status::${e}`" clearable/>
    </el-select>

    <el-select placeholder="Sort by" style="width: 150px;" v-model="sortValue">
      <el-option label="Created At ↓" value="createdAt-desc" />
      <el-option label="Created At ↑" value="createdAt-asc" />
      <el-option label="Updated At ↓" value="updatedAt-desc" />
      <el-option label="Updated At ↑" value="updatedAt-asc" />
      <el-option label="Firstname At ↑" value="firstname-asc" />
      <el-option label="Firstname At ↓" value="level-desc" />
      <el-option label="Lastname At ↑" value="lastname-asc" />
      <el-option label="Lastname At ↓" value="lastname-desc" />
    </el-select>

    <el-button style="padding:5px; height: 1.8rem" type="info" size="small" @click="handleReset"> Reset </el-button>
    <!--el-button type="info" @click="handleSearch"> Search </el-button-->
  </div>

  <div class="toolbar-right">
    <el-dropdown trigger="click">
      <el-icon :style="{ fontSize: '20px', margin: '2px' }" title="Columns" class="g-hover-icon"
      > <ArrowDown /> </el-icon>

      <template #dropdown>
        <el-dropdown-menu class="column-dropdown">
          <el-checkbox-group v-model="visibleColumns">
            <el-dropdown-item  class="no-hover" v-for="col in allColumns" :key="col.prop">
              <el-checkbox :value="col.prop"> {{ col.label }} </el-checkbox>
            </el-dropdown-item>
          </el-checkbox-group>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <el-button style="padding:5px" type="danger" @click="deleteSelected" :disabled="!selectedRows.length">
      Delete
    </el-button>



    <el-button style="padding:5px" type="primary" @click="openCreateAccount"> Create </el-button>
    <!--CreateAccount v-model:visible="openCreateAccount" @success="refresh" /-->
  </div>
</div>

<el-table
  style="margin-top: 10px;" empty-text="No accounts found" border
  v-loading="loading" :data="pageData.items" @selection-change="selectedRows = $event"
>
  <el-table-column type="selection" width="40" />

  <el-table-column label="Full Name" prop="fullName" width="120" sortable>
    <template #default="scope"> {{ scope.row.firstname }} {{ scope.row.lastname }} </template>
  </el-table-column>

  <!--el-table-column prop="id" label="ID" sortable :sort-method="(a, b) => a > b"/-->
  <el-table-column v-for="col in selectVisibleColumns"
    :prop="col.prop" :label="col.label" :key="col.prop"
    :sortable="col.sortable" :width="col.width" :formatter="col.formatter"
  />

  <el-table-column label="Actions" fixed="right" width="180">
    <template #default="scope">
      <div class="g-cell-center">
        <el-icon
          :style="{fontSize: '20px', margin: '2px'}" title="Copy Account" class="g-hover-icon"
          @click="copyToClipboard(scope.row)"
        > <CopyDocument /> </el-icon>

        <el-icon
          :style="{fontSize: '20px', margin: '2px'}" title="Edit Account" class="g-hover-icon"
          @click="editAccount(scope.row)"
        > <Edit /> </el-icon>

        <el-button
          size="small" style="padding:5px"
          :type="scope.row.status == 'activated' ? '' : 'warning'"
          @click="updateStatus(scope.row)"
        > {{ scope.row.status }} </el-button>
      </div>
    </template>
  </el-table-column>
</el-table>

<el-alert type="error" show-icon style="margin-top: 10px" v-if="error" :title="error"/>

<el-pagination
  class="pagination" layout="prev, pager, next, jumper, total, sizes"
  v-model:current-page="query.pageIndex" v-model:page-size="query.pageSize"
  :total="pageData.total" :page-sizes="[10, 20, 50]"
  @current-change="fetchData()"
/>
<!-- @size-change="" -->



<CreateAccount v-model:visible="createAccountVisible" @success="fetchData" />

<UpdateStatus
  v-model:visible="updateStatusVisible" :account="selectedAccount"
  @update:status="({ id, status }) => {
    const target = pageData.items.find(item => item.id === id)
    if (target) { target.status = status }
    selectedAccount.status = status; selectedAccount.newStatus = '';
  }"
/>

<EditAccount
  v-model:visible="editAccountVisible" :account="selectedAccount"
  @update:refresh="(form) => {
    const target = pageData.items.find(item => item.id === form.id)
    if (target) { Object.assign(target, form); delete target.password; }
  }"
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
