<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">内容审核</h1>
    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <n-data-table :columns="columns" :data="items" :loading="loading" :pagination="pagination" :bordered="false" :single-line="false" :row-key="(row: any) => row.id" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { get, put } from '@/api/request'
import { useMessage, NButton, NTag, type DataTableColumns } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])

const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '类型', key: 'target_type', width: 80, render: (row) => h(NTag, { size: 'small', bordered: false }, () => row.target_type || '-') },
  { title: '目标ID', key: 'target_id', width: 80 },
  { title: '状态', key: 'status', width: 90, render: (row) => h(NTag, { size: 'small', type: row.status === 0 ? 'warning' : row.status === 1 ? 'success' : 'error', bordered: false }, () => row.status === 0 ? '待审核' : row.status === 1 ? '已通过' : '已拒绝') },
  { title: '审核人ID', key: 'reviewer_id', width: 80, render: (row) => row.reviewer_id || '-' },
  { title: '原因', key: 'reason', ellipsis: { tooltip: true }, render: (row) => row.reason || '-' },
  { title: '时间', key: 'created_at', width: 160, render: (row) => row.created_at ? new Date(row.created_at).toLocaleString('zh-CN') : '' },
  { title: '操作', key: 'actions', width: 160, render: (row) => row.status === 0 ? h('div', { class: 'flex gap-2' }, [
    h(NButton, { size: 'small', type: 'success', quaternary: true, onClick: () => handleReview(row, 1) }, () => '通过'),
    h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleReview(row, 2) }, () => '拒绝')
  ]) : '-' }
]

const pagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (p: number) => { pagination.page = p; fetchItems() },
  onUpdatePageSize: (size: number) => { pagination.pageSize = size; pagination.page = 1; fetchItems() }
})

async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>('/admin/content-reviews', { page: pagination.page, page_size: pagination.pageSize })
    items.value = res.data?.items || []
    pagination.itemCount = res.data?.total || 0
  } catch { items.value = [] }
  finally { loading.value = false }
}

async function handleReview(row: any, status: number) {
  const reason = status === 1 ? '审核通过' : '审核拒绝'
  try {
    await put(`/admin/content-reviews/${row.id}`, { status, reason })
    row.status = status
    row.reason = reason
    message.success(status === 1 ? '已通过' : '已拒绝')
  } catch { message.error('操作失败') }
}

onMounted(fetchItems)
</script>
