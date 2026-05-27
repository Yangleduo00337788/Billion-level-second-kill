<template>
  <div>
    <h1 class="text-2xl font-bold text-dark mb-6">用户管理</h1>
    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <n-data-table
        :columns="columns"
        :data="users"
        :loading="loading"
        :pagination="paginationReactive"
        :bordered="false"
        :row-key="(row: any) => row.id"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { get, put } from '@/api/request'
import { useMessage, NButton, NTag, NSwitch, type DataTableColumns } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const users = ref<any[]>([])
const total = ref(0)
const page = ref(1)

const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 70 },
  { title: '用户名', key: 'username', ellipsis: { tooltip: true } },
  { title: '邮箱', key: 'email', ellipsis: { tooltip: true } },
  { title: '角色', key: 'role', width: 90, render: (row) => h(NTag, { size: 'small', type: row.role === 'admin' ? 'warning' : row.role === 'creator' ? 'info' : 'default', bordered: false }, () => row.role) },
  { title: '文章', key: 'article_count', width: 60 },
  { title: '粉丝', key: 'fans_count', width: 60 },
  { title: '状态', key: 'status', width: 80, render: (row) => h(NSwitch, { value: row.status === 1, 'onUpdate:value': (val: boolean) => toggleStatus(row, val) }) },
  { title: '注册时间', key: 'created_at', width: 160, render: (row) => formatDate(row.created_at) }
]

const paginationReactive = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (p: number) => { page.value = p; fetchUsers() },
  onUpdatePageSize: (size: number) => { paginationReactive.pageSize = size; page.value = 1; fetchUsers() }
})

function formatDate(d: string) {
  return d ? d.replace('T', ' ').substring(0, 19) : ''
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await get<any>('/admin/users', { page: page.value, page_size: paginationReactive.pageSize })
    users.value = res.data.items || []
    total.value = res.data.total
    paginationReactive.page = page.value
    paginationReactive.itemCount = res.data.total
  } catch { users.value = [] }
  finally { loading.value = false }
}

async function toggleStatus(row: any, val: boolean) {
  try {
    await put('/admin/users/' + row.id, { status: val ? 1 : 0 })
    row.status = val ? 1 : 0
    message.success(val ? '已启用' : '已禁用')
  } catch { message.error('操作失败') }
}

onMounted(fetchUsers)
</script>