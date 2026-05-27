<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">登录日志</h1>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" :pagination="pagination" remote />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, reactive, h } from 'vue'
import { NTag } from 'naive-ui'
import { get } from '@/api/request'
const loading = ref(false)
const items = ref<any[]>([])
const pagination = reactive({
  page: 1, pageSize: 20, pageCount: 1,
  onChange: (page: number) => { pagination.page = page; fetchItems() }
})
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '用户ID', key: 'user_id', width: 80 },
  { title: 'IP', key: 'ip', width: 130 },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => row.status === 1 ? '成功' : '失败' }) },
  { title: 'User Agent', key: 'user_agent', ellipsis: { tooltip: true } },
  { title: '时间', key: 'created_at', width: 170 }
]
async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>(`/admin/login-logs?page=${pagination.page}&page_size=${pagination.pageSize}`)
    const data = res.data
    items.value = data?.items || data?.data || (Array.isArray(data) ? data : [])
    pagination.pageCount = Math.ceil((data?.total || 0) / pagination.pageSize)
  } catch { items.value = [] } finally { loading.value = false }
}
onMounted(fetchItems)
</script>