<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-dark">公告管理</h1>
      <n-button type="primary" @click="showModal = true">
        <template #icon><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg></template>
        发布公告
      </n-button>
    </div>

    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />

    <n-modal v-model:show="showModal" preset="card" title="发布公告" style="max-width: 500px">
      <n-form ref="formRef" :model="form" label-placement="left" label-width="80">
        <n-form-item label="标题" path="title">
          <n-input v-model:value="form.title" placeholder="请输入标题" />
        </n-form-item>
        <n-form-item label="内容" path="content">
          <n-input v-model:value="form.content" type="textarea" placeholder="请输入内容" :rows="4" />
        </n-form-item>
        <n-form-item label="优先级" path="priority">
          <n-select v-model:value="form.priority" :options="priorityOptions" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="submitting" @click="handleSubmit">确定</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, NTag, NSpace, useMessage } from 'naive-ui'
import { get, post, put, del } from '@/api/request'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const items = ref<any[]>([])
const form = ref({ title: '', content: '', priority: 0 })

const priorityOptions = [
  { label: '普通', value: 0 },
  { label: '重要', value: 1 },
  { label: '紧急', value: 2 }
]

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  { title: '内容', key: 'content', ellipsis: { tooltip: true } },
  {
    title: '优先级', key: 'priority', width: 100,
    render: (row: any) => h(NTag, { type: ['', 'warning', 'error'][row.priority] || 'default', size: 'small' }, { default: () => ['普通', '重要', '紧急'][row.priority] || '普通' })
  },
  {
    title: '状态', key: 'status', width: 80,
    render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small' }, { default: () => row.status === 1 ? '发布' : '草稿' })
  },
  { title: '时间', key: 'created_at', width: 180 },
  {
    title: '操作', key: 'actions', width: 120,
    render: (row: any) => h(NSpace, { size: 'small' }, {
      default: () => [
        h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => toggleStatus(row) }, { default: () => row.status === 1 ? '下线' : '发布' }),
        h(NButton, { size: 'small', quaternary: true, type: 'error', onClick: () => handleDelete(row.id) }, { default: () => '删除' })
      ]
    })
  }
]

async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>('/admin/announcements')
    items.value = res.data?.items || []
  } catch {} finally { loading.value = false }
}

async function handleSubmit() {
  if (!form.value.title || !form.value.content) {
    message.warning('请填写完整')
    return
  }
  submitting.value = true
  try {
    await post('/admin/announcements', form.value)
    message.success('发布成功')
    showModal.value = false
    form.value = { title: '', content: '', priority: 0 }
    fetchItems()
  } catch {} finally { submitting.value = false }
}

async function toggleStatus(row: any) {
  try {
    await put(`/admin/announcements/${row.id}`, { status: row.status === 1 ? 0 : 1 })
    message.success('操作成功')
    fetchItems()
  } catch {}
}

async function handleDelete(id: number) {
  try {
    await del(`/admin/announcements/${id}`)
    message.success('删除成功')
    fetchItems()
  } catch {}
}

onMounted(fetchItems)
</script>