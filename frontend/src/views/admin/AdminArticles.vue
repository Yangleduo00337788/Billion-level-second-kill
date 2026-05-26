<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">文章管理</h1>
    <n-data-table :columns="columns" :data="articles" :loading="loading" :pagination="pagination" :bordered="false" class="bg-white rounded-xl" />
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { get, del } from '@/api/request'
import { useMessage, NButton, NTag, useDialog } from 'naive-ui'

const message = useMessage()
const dialog = useDialog()
const loading = ref(false)
const articles = ref<any[]>([])
const total = ref(0)
const page = ref(1)

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { size: 'small', type: row.status === 'published' ? 'success' : 'default' }, () => row.status) },
  { title: '阅读', key: 'view_count', width: 60 },
  { title: '点赞', key: 'like_count', width: 60 },
  { title: 'AI', key: 'is_ai', width: 50, render: (row: any) => row.is_ai ? '🤖' : '' },
  { title: '时间', key: 'created_at', width: 120 },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NButton, { size: 'tiny', type: 'error', quaternary: true, onClick: () => handleDelete(row) }, () => '删除') }
]

const pagination = { page, pageSize: 20, itemCount: total, onChange: (p: number) => { page.value = p; fetchArticles() } }

async function fetchArticles() {
  loading.value = true
  try {
    const res = await get<any>('/admin/articles', { page: page.value, page_size: 20 })
    articles.value = res.data.items || []
    total.value = res.data.total
  } catch { articles.value = [] }
  finally { loading.value = false }
}

function handleDelete(row: any) {
  dialog.warning({ title: '确认删除', content: `确定删除文章「${row.title}」？`, positiveText: '删除', negativeText: '取消', onPositiveClick: async () => {
    try { await del(`/admin/articles/${row.id}`); articles.value = articles.value.filter((a: any) => a.id !== row.id); message.success('已删除') }
    catch { message.error('删除失败') }
  }})
}

onMounted(fetchArticles)
</script>
