<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">举报管理</h1>

    <n-tabs v-model:value="activeTab" @update:value="fetchItems">
      <n-tab-pane name="pending" tab="待处理" />
      <n-tab-pane name="processed" tab="已处理" />
      <n-tab-pane name="ignored" tab="已忽略" />
    </n-tabs>

    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" class="mt-4" />

    <n-modal v-model:show="showModal" preset="card" title="处理举报" style="max-width: 500px">
      <n-form label-placement="left" label-width="80">
        <n-form-item label="举报原因">
          <n-input :value="currentItem?.reason" type="textarea" :rows="3" readonly />
        </n-form-item>
        <n-form-item label="处理结果">
          <n-input v-model:value="handleResult" type="textarea" placeholder="请输入处理结果" :rows="3" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="handleIgnore">忽略</n-button>
        <n-button type="primary" :loading="submitting" @click="handleApprove">确认处理</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, NTag, useMessage } from 'naive-ui'
import { get, put } from '@/api/request'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const activeTab = ref('pending')
const items = ref<any[]>([])
const currentItem = ref<any>(null)
const handleResult = ref('')

const statusLabels = ['待处理', '已处理', '已忽略']
const statusTypes = ['warning', 'success', 'default']

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '类型', key: 'target_type', width: 80, render: (row: any) => h(NTag, { size: 'small' }, { default: () => row.target_type }) },
  { title: '目标ID', key: 'target_id', width: 80 },
  { title: '举报原因', key: 'reason', ellipsis: { tooltip: true } },
  { title: '举报人ID', key: 'reporter_id', width: 100 },
  {
    title: '状态', key: 'status', width: 100,
    render: (row: any) => h(NTag, { type: (statusTypes[row.status] || 'default') as any, size: 'small' }, { default: () => statusLabels[row.status] || '未知' })
  },
  { title: '处理结果', key: 'result', ellipsis: { tooltip: true } },
  { title: '时间', key: 'created_at', width: 180 },
  {
    title: '操作', key: 'actions', width: 100,
    render: (row: any) => {
      if (row.status !== 0) return null
      return h(NButton, { size: 'small', type: 'primary', quaternary: true, onClick: () => openHandle(row) }, { default: () => '处理' })
    }
  }
]

async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>('/admin/reports')
    const allItems = res.data?.items || []
    const statusMap: Record<string, number> = { pending: 0, processed: 1, ignored: 2 }
    items.value = allItems.filter((i: any) => i.status === statusMap[activeTab.value])
  } catch {} finally { loading.value = false }
}

function openHandle(row: any) {
  currentItem.value = row
  handleResult.value = ''
  showModal.value = true
}

async function handleApprove() {
  if (!handleResult.value) { message.warning('请输入处理结果'); return }
  submitting.value = true
  try {
    await put(`/admin/reports/${currentItem.value.id}`, { status: 1, result: handleResult.value })
    message.success('处理成功')
    showModal.value = false
    fetchItems()
  } catch {} finally { submitting.value = false }
}

async function handleIgnore() {
  try {
    await put(`/admin/reports/${currentItem.value.id}`, { status: 2, result: '已忽略' })
    message.success('已忽略')
    showModal.value = false
    fetchItems()
  } catch {}
}

onMounted(fetchItems)
</script>