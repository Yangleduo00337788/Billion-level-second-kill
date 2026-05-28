<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">邀请码管理</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4 flex gap-3">
      <n-button type="primary" @click="handleCreate">生成邀请码</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, NTag, useMessage } from 'naive-ui'
import { get, post } from '@/api/request'
const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '邀请码', key: 'code', width: 150, render: (row: any) => h('span', { class: 'font-mono text-sm' }, row.code) },
  { title: '最大使用', key: 'max_uses', width: 100 },
  { title: '已使用', key: 'used_count', width: 80 },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: (row.used_count || 0) < row.max_uses ? 'success' : 'default', size: 'small' }, { default: () => (row.used_count || 0) < row.max_uses ? '可用' : '已用完' }) },
  { title: '过期时间', key: 'expires_at', width: 170, render: (row: any) => row.expires_at ? new Date(row.expires_at).toLocaleString('zh-CN') : '永不过期' },
  { title: '创建时间', key: 'created_at', width: 170, render: (row: any) => row.created_at ? new Date(row.created_at).toLocaleString('zh-CN') : '-' }
]
async function fetchItems() {
  loading.value = true
  try { const res = await get<any>('/admin/invite-codes'); items.value = Array.isArray(res.data) ? res.data : [] } catch { items.value = [] } finally { loading.value = false }
}
async function handleCreate() {
  try { await post('/admin/invite-codes', {}); message.success('邀请码已生成'); fetchItems() } catch {}
}
onMounted(fetchItems)
</script>