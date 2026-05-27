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
const columns = [
  { title: 'ID', key: 'id', width: 60, fixed: 'left' as const },
  { title: '用户', key: 'username', width: 100, fixed: 'left' as const },
  {
    title: '操作', key: 'action', width: 120,
    render: (row: any) => {
      const action = row.action || ''
      let type = 'default'
      if (action.includes('POST') || action.includes('创建')) type = 'success'
      else if (action.includes('PUT') || action.includes('编辑') || action.includes('更新')) type = 'info'
      else if (action.includes('DELETE') || action.includes('删除')) type = 'error'
      else if (action.includes('登录')) type = 'warning'
      return h(NTag, { type: type as any, size: 'small' }, { default: () => action })
    }
  },
  { title: '目标', key: 'target', width: 150, ellipsis: { tooltip: true } },
  { title: '详情', key: 'detail', width: 200, ellipsis: { tooltip: true } },
  { title: 'IP', key: 'ip', width: 130 },
  { title: '时间', key: 'created_at', width: 170 }
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