<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">用户管理</h1>
    <n-data-table :columns="columns" :data="users" :loading="loading" :pagination="pagination" :bordered="false" class="bg-white rounded-xl" />
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { get, put } from '@/api/request'
import { useMessage, NButton, NTag, NSwitch } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const users = ref<any[]>([])
const total = ref(0)
const page = ref(1)

const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '用户名', key: 'username' },
  { title: '邮箱', key: 'email' },
  { title: '角色', key: 'role', render: (row: any) => h(NTag, { size: 'small', type: row.role === 'admin' ? 'warning' : 'info' }, () => row.role) },
  { title: '文章数', key: 'article_count', width: 80 },
  { title: '状态', key: 'status', render: (row: any) => h(NSwitch, { value: row.status === 1, 'onUpdate:value': (val: boolean) => toggleStatus(row, val) }) },
  { title: '注册时间', key: 'created_at' }
]

const pagination = { page: page, pageSize: 20, itemCount: total, onChange: (p: number) => { page.value = p; fetchUsers() } }

async function fetchUsers() {
  loading.value = true
  try {
    const res = await get<any>('/admin/users', { page: page.value, page_size: 20 })
    users.value = res.data.items || []
    total.value = res.data.total
  } catch { users.value = [] }
  finally { loading.value = false }
}

async function toggleStatus(row: any, val: boolean) {
  try {
    await put(`/admin/users/${row.id}`, { status: val ? 1 : 0 })
    row.status = val ? 1 : 0
    message.success(val ? '已启用' : '已禁用')
  } catch { message.error('操作失败') }
}

onMounted(fetchUsers)
</script>
