<template>
  <div>
    <h1 class="text-2xl font-bold text-dark mb-6">文章管理</h1>
    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <n-data-table
        :columns="columns"
        :data="articles"
        :loading="loading"
        :pagination="paginationReactive"
        :bordered="false"
        :row-key="(row: any) => row.id"
      />
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
const articles = ref<any[]>([])
const total = ref(0)
const page = ref(1)

const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  { title: '作者ID', key: 'user_id', width: 70 },
  { title: '状态', key: 'status', width: 80, render: (row) => h(NTag, { size: 'small', type: row.status === 'published' ? 'success' : 'default', bordered: false }, () => row.status === 'published' ? '已发布' : '草稿') },
  { title: '浏览', key: 'view_count', width: 60 },
  { title: '点赞', key: 'like_count', width: 60 },
  { title: '评论', key: 'comment_count', width: 60 },
  { title: 'AI', key: 'is_ai', width: 50, render: (row) => row.is_ai ? h('span', { class: 'text-xs text-orange-500 font-medium' }, 'AI') : '' },
  { title: '时间', key: 'created_at', width: 160, render: (row) => formatDate(row.created_at) },
  { title: '操作', key: 'actions', width: 80, render: (row) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row) }, () => '删除') }
]

const paginationReactive = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (p: number) => { page.value = p; fetchArticles() },
  onUpdatePageSize: (size: number) => { paginationReactive.pageSize = size; page.value = 1; fetchArticles() }
})

function formatDate(d: string) {
  return d ? d.replace('T', ' ').substring(0, 19) : ''
}

async function fetchArticles() {
  loading.value = true
  try {
    const res = await get<any>('/admin/articles', { page: page.value, page_size: paginationReactive.pageSize })
    articles.value = res.data.items || []
    total.value = res.data.total
    paginationReactive.page = page.value
    paginationReactive.itemCount = res.data.total
  } catch { articles.value = [] }
  finally { loading.value = false }
}

function handleDelete(row: any) {
  dialog.warning({
    title: '确认删除',
    content: '确定删除文章「' + row.title + '」？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await del('/admin/articles/' + row.id)
        articles.value = articles.value.filter((a: any) => a.id !== row.id)
        message.success('已删除')
      } catch { message.error('删除失败') }
    }
  })
}

onMounted(fetchArticles)
</script>