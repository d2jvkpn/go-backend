<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

// import { hello } from "@/js/_hello.js"
// hello()
import { service } from "@/js/utils/request.js";

//
const allColumns = [
  { prop: 'id',         label: 'ID' },
  { prop: 'firstname',  label: 'Firstname' },
  { prop: 'lastname',   label: 'Lastname' },
  { prop: 'email',      label: 'Email' },
  { prop: 'level',      label: 'Level' },
  { prop: 'labels',     label: 'Labels' },
  { prop: 'status',     label: 'Status', sortable: true },
  { prop: 'createdAt',  label: 'CreatedAt', sortable: true },
  { prop: 'updatedAt', label: 'updatedAt', sortable: true },
]

const visibleColumns = ref(['firstname', 'lastname', 'email', 'level', 'labels', 'status', 'createdAt'])

const selectVisibleColumns = computed(() =>
  allColumns.filter(col => visibleColumns.value.includes(col.prop))
)

const pageData = ref({ total: 0, items: [] })
const loading = ref(false)
const error = ref(null);

const fetchData = async () => {
  loading.value = true;
  error.value = null;

  try {
     const data = await service.get("/api/v1/auth/account/query_accounts", {}, { params: query.value });

     if (query.value.pageIndex == 1) {
       pageData.value.total = data.total;
     }

     pageData.value.items = data.items;
  } catch (err) {
    error.value = err.response?.data?.msg || err.msg;
  } finally {
    loading.value = false;
  }
};

//
const levels = ['admin', 'editor', 'reviewer', 'user', 'guest'];
const statuses = ["created", "activated", "blocked"]

const query = ref({ pageIndex: 1, pageSize: 10, keyword: '', level: '', status: 'activated' })

const handleReset = async () => {
  query.value.pageIndex = 1;
  query.value.keyword = '';
  query.value.level = '';
  query.value.status = 'activated';
  await fetchData();
}

const handleSearch = async () => {
  query.value.pageIndex = 1;
  await fetchData();
};

const hanldePageIndexChange = async (v) => {
  console.log(`==> hanldePageIndexChange: ${v}`)
  query.value.pageIndex = v;
  await fetchData();;
}

const handlePageSizeChange = async (v) => {
  console.log(`==> handlePageSizeChange: ${v}`)
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

  /*
  ElMessageBox.confirm(
    `Are you sure you want to delete ${idsToDelete.length} account${s}?`,
    `Delete Account${s} Confirmation`,
    { type: 'warning', confirmButtonText: 'Yes', cancelButtonText: 'No' }
  )
  .then(() => {
    // TODO
    selectedRows.value = []
    ElMessage.success(`Deleted ${idsToDelete.length} account${s}`)
  })
  .catch(() => {
    ElMessage.info(`Delete account${s} canceled.`)
  })
  */

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

/*
const keyword = ref('');

watch(keyword, (newVal, oldVal) => {
  console.log(`keyword changed from ${oldVal} to ${newVal}`);
});
*/

/*
import debounce from 'lodash/debounce'

const debouncedFetch = debounce(async () => {
  console.log('🔍 用户停止输入，开始搜索')
  await fetchData()
}, 500)  // 500 毫秒内没输入，才触发

watch(() => query.value.keyword, () => {
  debouncedFetch()
})
*/

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

    <el-select v-model="query.level" placeholder="level" style="width: 150px" clearable>
      <el-option v-for="e in levels" :value="e" :label="e" :key="`account::level::${e}`" />
    </el-select>

    <el-select v-model="query.status" placeholder="status" style="width: 150px" clearable>
      <el-option v-for="e in statuses" :value="e" :label="e" :key="`account::status::${e}`" />
    </el-select>

    <el-button @click="handleReset"> Reset </el-button>
    <!--el-button @click="handleSearch"> Search </el-button-->
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

<el-table
  :data="pageData.items"
  style="margin-top: 10px;"
  @selection-change="selectedRows = $event"
  border
  v-loading="loading"
  empty-text="No accounts found"
>

  <el-table-column type="selection" width="50" />

  <!--el-table-column prop="id" label="ID" sortable :sort-method="customSort"/-->

  <el-table-column v-for="col in selectVisibleColumns"
    :prop="col.prop" :label="col.label" :key="col.prop" :sortable="col.sortable"
  />
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
