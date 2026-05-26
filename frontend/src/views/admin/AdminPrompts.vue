<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">Prompt 管理</h1>
    <n-data-table :columns="columns" :data="prompts" :loading="loading" :pagination="pagination" :bordered="false" class="bg-white rounded-xl" />
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { get, del } from '@/api/request'
import { useMessage, NButton, NTag, useDialog } from 'naive-ui'

const message = useMessage()
const dialog = useDialog()
const loading = ref(false)
const prompts = ref<any[]>([])
const total = ref(0)
const page = ref(1)

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '标题', key: 'title' },
  { title: '分类', key: 'category', width: 80, render: (row: any) => h(NTag, { size: 'small' }, () => row.category) },
  { title: '使用次数', key: 'usage_count', width: 80 },
  { title: '状态', key: 'status', width: 60, render: (row: any) => h(NTag, { size: 'small', type: row.status === 1 ? 'success' : 'default' }, () => row.status === 1 ? '正常' : '禁用') },
  { title: '时间', key: 'created_at', width: 120 },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NButton, { size: 'tiny', type: 'error', quaternary: true, onClick: () => handleDelete(row) }, () => '删除') }
]

const pagination = { page, pageSize: 20, itemCount: total, onChange: (p: number) => { page.value = p; fetchPrompts() } }

async function fetchPrompts() {
  loading.value = true
  try {
    const res = await get<any>('/admin/prompts', { page: page.value, page_size: 20 })
    prompts.value = res.data.items || []
    total.value = res.data.total
  } catch { prompts.value = [] }
  finally { loading.value = false }
}

function handleDelete(row: any) {
  dialog.warning({ title: '确认删除', content: `确定删除 Prompt「${row.title}」？`, positiveText: '删除', negativeText: '取消', onPositiveClick: async () => {
    try { await del(`/admin/prompts/${row.id}`); prompts.value = prompts.value.filter((p: any) => p.id !== row.id); message.success('已删除') }
    catch { message.error('删除失败') }
  }})
}

onMounted(fetchPrompts)
</script>
