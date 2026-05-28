<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">操作日志</h1>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" :pagination="pagination" remote />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, reactive, h } from 'vue'
import { NTag, NTooltip } from 'naive-ui'
import { get } from '@/api/request'

const loading = ref(false)
const items = ref<any[]>([])
const pagination = reactive({
  page: 1, pageSize: 20, pageCount: 1,
  showSizePicker: true, pageSizes: [20, 50, 100],
  onChange: (page: number) => { pagination.page = page; fetchItems() },
  onUpdatePageSize: (size: number) => { pagination.pageSize = size; pagination.page = 1; fetchItems() }
})

function getActionType(action: string): string {
  if (!action) return 'default'
  if (action.includes('注册') || action.includes('创建') || action.includes('发布')) return 'success'
  if (action.includes('编辑') || action.includes('更新') || action.includes('修改')) return 'info'
  if (action.includes('删除') || action.includes('封禁')) return 'error'
  if (action.includes('登录') || action.includes('登出')) return 'warning'
  if (action.includes('关注') || action.includes('取消')) return 'primary'
  if (action.includes('点赞') || action.includes('收藏')) return 'error'
  return 'default'
}

function formatDetail(row: any): string {
  const parts = []
  if (row.username) parts.push(`用户「${row.username}」`)
  if (row.action) parts.push(row.action)
  if (row.target) parts.push(`于「${row.target}」`)
  if (row.detail) parts.push(`— ${row.detail}`)
  return parts.join(' ') || '-'
}

const columns = [
  { title: 'ID', key: 'id', width: 60, fixed: 'left' as const },
  {
    title: '操作用户', key: 'username', width: 100, fixed: 'left' as const,
    render: (row: any) => h('span', { class: 'text-sm font-medium' }, row.username || `用户#${row.user_id}` || '-')
  },
  {
    title: '操作类型', key: 'action', width: 130,
    render: (row: any) => h(NTag, { type: getActionType(row.action) as any, size: 'small', round: true }, { default: () => row.action || '-' })
  },
  {
    title: '操作详情', key: 'detail', minWidth: 250,
    render: (row: any) => {
      const text = formatDetail(row)
      return h(NTooltip, { trigger: 'hover' }, {
        trigger: () => h('span', { class: 'text-sm text-gray-600 truncate block max-w-md' }, text),
        default: () => text
      })
    }
  },
  {
    title: 'IP 地址', key: 'ip', width: 140,
    render: (row: any) => h('span', { class: 'text-xs font-mono text-gray-500' }, row.ip || '-')
  },
  {
    title: '操作时间', key: 'created_at', width: 170,
    render: (row: any) => h('span', { class: 'text-xs text-gray-400' }, row.created_at ? new Date(row.created_at).toLocaleString('zh-CN') : '-')
  }
]

async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>(`/admin/audit-logs?page=${pagination.page}&page_size=${pagination.pageSize}`)
    const data = res.data
    items.value = data?.items || data?.data || (Array.isArray(data) ? data : [])
    pagination.pageCount = Math.ceil((data?.total || 0) / pagination.pageSize)
  } catch { items.value = [] } finally { loading.value = false }
}
onMounted(fetchItems)
</script>