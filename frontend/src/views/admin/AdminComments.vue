<template>
  <div>
    <h1 class="text-2xl font-bold text-dark mb-6">评论管理</h1>
    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <n-data-table :columns="columns" :data="comments" :loading="loading" :pagination="paginationReactive" :bordered="false" :row-key="(row: any) => row.id" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { get, del } from '@/api/request'
import { useMessage, NButton, useDialog, type DataTableColumns } from 'naive-ui'

const message = useMessage()
const dialog = useDialog()
const loading = ref(false)
const comments = ref<any[]>([])
const page = ref(1)

const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '内容', key: 'content', ellipsis: { tooltip: true } },
  { title: '用户ID', key: 'user_id', width: 70 },
  { title: '文章ID', key: 'article_id', width: 70 },
  { title: '时间', key: 'created_at', width: 160, render: (row) => row.created_at ? row.created_at.replace('T', ' ').substring(0, 19) : '' },
  { title: '操作', key: 'actions', width: 80, render: (row) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row) }, () => '删除') }
]

const paginationReactive = reactive({ page: 1, pageSize: 20, itemCount: 0, onChange: (p: number) => { page.value = p; fetchComments() } })

async function fetchComments() {
  loading.value = true
  try {
    const res = await get<any>('/admin/comments', { page: page.value, page_size: 20 })
    comments.value = res.data?.items || []
    paginationReactive.page = page.value
    paginationReactive.itemCount = res.data?.total || 0
  } catch { comments.value = [] }
  finally { loading.value = false }
}

function handleDelete(row: any) {
  dialog.warning({ title: '确认删除', content: '确定删除该评论？', positiveText: '删除', negativeText: '取消',
    onPositiveClick: async () => {
      try { await del('/admin/comments/' + row.id); comments.value = comments.value.filter((c: any) => c.id !== row.id); message.success('已删除') }
      catch { message.error('删除失败') }
    }
  })
}

onMounted(fetchComments)
</script>