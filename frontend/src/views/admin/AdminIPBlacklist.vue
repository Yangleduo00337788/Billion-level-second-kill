<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">IP 黑名单</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <div class="flex gap-3">
        <n-input v-model:value="searchIP" placeholder="输入 IP 地址" clearable style="width: 200px" />
        <n-input v-model:value="addReason" placeholder="封禁原因" clearable style="width: 200px" />
        <n-button type="primary" @click="handleAdd">添加 IP</n-button>
      </div>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, useMessage } from 'naive-ui'
import { get, post, del } from '@/api/request'
const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])
const searchIP = ref('')
const addReason = ref('')
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: 'IP 地址', key: 'ip', width: 150 },
  { title: '原因', key: 'reason', ellipsis: { tooltip: true } },
  { title: '添加时间', key: 'created_at', width: 170 },
  { title: '操作', key: 'actions', width: 80, render: (row: any) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row.id) }, { default: () => '删除' }) }
]
async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>('/admin/ip-blacklist')
    items.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
  } catch { items.value = [] } finally { loading.value = false }
}
async function handleAdd() {
  if (!searchIP.value) { message.warning('请输入 IP 地址'); return }
  try {
    await post('/admin/ip-blacklist', { ip: searchIP.value, reason: addReason.value })
    message.success('添加成功')
    searchIP.value = ''
    addReason.value = ''
    fetchItems()
  } catch {}
}
async function handleDelete(id: number) {
  try { await del(`/admin/ip-blacklist/${id}`); message.success('删除成功'); fetchItems() } catch {}
}
onMounted(fetchItems)
</script>