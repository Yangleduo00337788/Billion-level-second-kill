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
  showSizePicker: true, pageSizes: [20, 50, 100],
  onChange: (p: number) => { pagination.page = p; fetchItems() },
  onUpdatePageSize: (size: number) => { pagination.pageSize = size; pagination.page = 1; fetchItems() }
})

function getTypeTag(type: string) {
  const map: Record<string, { label: string; color: string }> = {
    like: { label: '点赞', color: 'error' },
    comment: { label: '评论', color: 'info' },
    follow: { label: '关注', color: 'success' }
  }
  return map[type] || { label: type || '未知', color: 'default' }
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '接收用户', key: 'username', width: 100,
    render: (row: any) => h('span', { class: 'text-sm' }, row.username || `#${row.user_id}`)
  },
  {
    title: '触发用户', key: 'actor_name', width: 100,
    render: (row: any) => h('span', { class: 'text-sm' }, row.actor_name || `#${row.actor_id}`)
  },
  {
    title: '类型', key: 'type', width: 80,
    render: (row: any) => {
      const t = getTypeTag(row.type)
      return h(NTag, { type: t.color as any, size: 'small', round: true }, { default: () => t.label })
    }
  },
  {
    title: '内容', key: 'content', minWidth: 200, ellipsis: { tooltip: true },
    render: (row: any) => h('span', { class: 'text-sm text-gray-600' }, row.content || '-')
  },
  {
    title: '已读', key: 'is_read', width: 80,
    render: (row: any) => h(NTag, { type: row.is_read ? 'success' : 'warning', size: 'small' }, { default: () => row.is_read ? '已读' : '未读' })
  },
  {
    title: '时间', key: 'created_at', width: 170,
    render: (row: any) => h('span', { class: 'text-xs text-gray-400' }, row.created_at ? new Date(row.created_at).toLocaleString('zh-CN') : '-')
  }
]

async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>(`/admin/notifications?page=${pagination.page}&page_size=${pagination.pageSize}`)
    const data = res.data
    items.value = data?.items || data?.data || (Array.isArray(data) ? data : [])
    pagination.pageCount = Math.ceil((data?.total || 0) / pagination.pageSize)
  } catch { items.value = [] } finally { loading.value = false }
}
onMounted(fetchItems)
</script>