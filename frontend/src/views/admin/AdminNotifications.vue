<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">通知管理</h1>
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
  onChange: (p: number) => { pagination.page = p; fetchItems() }
})
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '接收用户ID', key: 'user_id', width: 100 },
  { title: '触发用户ID', key: 'actor_id', width: 100 },
  {
    title: '类型', key: 'type', width: 80,
    render: (row: any) => {
      const map: Record<string, any> = { like: { label: '点赞', type: 'error' }, comment: { label: '评论', type: 'info' }, follow: { label: '关注', type: 'success' } }
      const t = map[row.type] || { label: row.type, type: 'default' }
      return h(NTag, { type: t.type, size: 'small' }, { default: () => t.label })
    }
  },
  { title: '内容', key: 'content', ellipsis: { tooltip: true } },
  { title: '目标ID', key: 'target_id', width: 80 },
  { title: '状态', key: 'is_read', width: 80, render: (row: any) => h(NTag, { type: row.is_read ? 'success' : 'warning', size: 'small' }, { default: () => row.is_read ? '已读' : '未读' }) },
  { title: '时间', key: 'created_at', width: 170 }
]
async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>(`/notifications?page=${pagination.page}&page_size=${pagination.pageSize}`)
    const data = res.data
    items.value = data?.items || data?.data || (Array.isArray(data) ? data : [])
    pagination.pageCount = Math.ceil((data?.total || 0) / pagination.pageSize)
  } catch { items.value = [] } finally { loading.value = false }
}
onMounted(fetchItems)
</script>