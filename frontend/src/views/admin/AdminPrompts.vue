<template>
  <div>
    <h1 class="text-2xl font-bold text-dark mb-6">Prompt 管理</h1>
    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <n-data-table :columns="columns" :data="prompts" :loading="loading" :pagination="paginationReactive" :bordered="false" :row-key="(row: any) => row.id" />
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { get, del } from '@/api/request'
import { useMessage, NButton, NTag, useDialog, type DataTableColumns } from 'naive-ui'
const message = useMessage()
const dialog = useDialog()
const loading = ref(false)
const prompts = ref<any[]>([])
const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  { title: '作者ID', key: 'user_id', width: 70 },
  { title: '分类', key: 'category', width: 80, render: (row) => h(NTag, { size: 'small', bordered: false }, () => row.category || '-') },
  { title: '模型', key: 'model', width: 90 },
  { title: '使用', key: 'usage_count', width: 60 },
  { title: '点赞', key: 'like_count', width: 60 },
  { title: '状态', key: 'status', width: 70, render: (row) => h(NTag, { size: 'small', type: row.status === 1 ? 'success' : 'error', bordered: false }, () => row.status === 1 ? '正常' : '禁用') },
  { title: '时间', key: 'created_at', width: 160, render: (row) => row.created_at ? new Date(row.created_at).toLocaleString('zh-CN') : '' },
  { title: '操作', key: 'actions', width: 80, render: (row) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row) }, () => '删除') }
]
const paginationReactive = reactive({ page: 1, pageSize: 20, itemCount: 0, showSizePicker: true, pageSizes: [10, 20, 50], onChange: (p: number) => { paginationReactive.page = p; fetchPrompts() }, onUpdatePageSize: (size: number) => { paginationReactive.pageSize = size; paginationReactive.page = 1; fetchPrompts() } })
async function fetchPrompts() {
  loading.value = true
  try { const res = await get<any>('/admin/prompts', { page: paginationReactive.page, page_size: paginationReactive.pageSize }); prompts.value = res.data.items || []; paginationReactive.itemCount = res.data.total || 0 } catch { prompts.value = [] } finally { loading.value = false }
}
function handleDelete(row: any) {
  dialog.warning({
    title: '确认删除',
    content: '确定删除 Prompt「' + row.title + '」？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try { await del('/admin/prompts/' + row.id); prompts.value = prompts.value.filter((p: any) => p.id !== row.id); message.success('已删除') } catch { message.error('删除失败') }
    }
  })
}
onMounted(fetchPrompts)
</script>