<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">IP 黑名单</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <div class="flex gap-3 items-end">
        <div class="flex-1">
          <label class="text-xs text-gray-500 mb-1 block">选择用户</label>
          <n-select
            v-model:value="selectedUser"
            :options="userOptions"
            :loading="searchLoading"
            filterable
            remote
            :remote-method="searchUsers"
            placeholder="搜索用户名或邮箱"
            clearable
            style="width: 100%"
          />
        </div>
        <div style="width: 200px">
          <label class="text-xs text-gray-500 mb-1 block">封禁原因</label>
          <n-select v-model:value="addReason" :options="reasonOptions" placeholder="选择或输入原因" filterable tag clearable style="width: 100%" />
        </div>
        <n-button type="error" @click="handleAdd" :disabled="!selectedUser">
          封禁用户 IP
        </n-button>
      </div>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, NTag, useMessage } from 'naive-ui'
import { get, post, del } from '@/api/request'

const message = useMessage()
const loading = ref(false)
const searchLoading = ref(false)
const items = ref<any[]>([])
const selectedUser = ref<number | null>(null)
const addReason = ref('')
const userOptions = ref<any[]>([])

const reasonOptions = [
  { label: '恶意刷接口', value: '恶意刷接口' },
  { label: '发布违规内容', value: '发布违规内容' },
  { label: '垃圾广告', value: '垃圾广告' },
  { label: '恶意注册', value: '恶意注册' },
  { label: '攻击行为', value: '攻击行为' },
  { label: '其他', value: '其他' }
]

async function loadAllUsers() {
  searchLoading.value = true
  try {
    const res = await get<any>('/admin/users?page=1&page_size=100')
    const users = res.data?.items || []
    userOptions.value = users.map((u: any) => ({
      label: `${u.username} (${u.email})${u.status === 0 ? ' [已封禁]' : ''}`,
      value: u.id
    }))
  } catch {} finally { searchLoading.value = false }
}

async function searchUsers(query: string) {
  if (!query) { loadAllUsers(); return }
  searchLoading.value = true
  try {
    const res = await get<any>(`/admin/users/search?q=${encodeURIComponent(query)}`)
    const users = Array.isArray(res.data) ? res.data : []
    userOptions.value = users.map((u: any) => ({
      label: `${u.username} (${u.email})${u.status === 0 ? ' [已封禁]' : ''}`,
      value: u.id
    }))
  } catch {} finally { searchLoading.value = false }
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: 'IP 地址', key: 'ip', width: 150,
    render: (row: any) => h('span', { class: 'font-mono text-sm' }, row.ip || '-')
  },
  {
    title: '关联用户', key: 'username', width: 120,
    render: (row: any) => h('span', { class: 'text-sm' }, row.username || (row.user_id ? `#${row.user_id}` : '-'))
  },
  {
    title: '封禁原因', key: 'reason', minWidth: 150, ellipsis: { tooltip: true },
    render: (row: any) => h('span', { class: 'text-sm text-gray-600' }, row.reason || '-')
  },
  {
    title: '添加时间', key: 'created_at', width: 170,
    render: (row: any) => h('span', { class: 'text-xs text-gray-400' }, row.created_at ? new Date(row.created_at).toLocaleString('zh-CN') : '-')
  },
  {
    title: '操作', key: 'actions', width: 80,
    render: (row: any) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row.id) }, { default: () => '解封' })
  }
]

async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>('/admin/ip-blacklist')
    items.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
  } catch { items.value = [] } finally { loading.value = false }
}

async function handleAdd() {
  if (!selectedUser.value) { message.warning('请选择用户'); return }
  try {
    await post('/admin/ip-blacklist', {
      user_id: selectedUser.value,
      reason: addReason.value || '管理员封禁'
    })
    message.success('封禁成功')
    selectedUser.value = null
    addReason.value = ''
    fetchItems()
  } catch {}
}

async function handleDelete(id: number) {
  try { await del(`/admin/ip-blacklist/${id}`); message.success('解封成功'); fetchItems() } catch {}
}

onMounted(() => {
  fetchItems()
  loadAllUsers()
})
</script>