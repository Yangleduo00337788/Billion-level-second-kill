<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">系统配置</h1>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NInput, NButton, useMessage } from 'naive-ui'
import { get, put } from '@/api/request'
const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])
const editingId = ref<number | null>(null)
const editingValue = ref('')
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: 'Key', key: 'key', width: 200 },
  {
    title: 'Value', key: 'value', minWidth: 300,
    render: (row: any) => {
      if (editingId.value === row.id) {
        return h(NInput, { value: editingValue.value, onUpdateValue: (v: string) => { editingValue.value = v }, size: 'small' })
      }
      return h('span', { class: 'text-sm' }, row.value)
    }
  },
  { title: '说明', key: 'desc', ellipsis: { tooltip: true } },
  {
    title: '操作', key: 'actions', width: 120,
    render: (row: any) => {
      if (editingId.value === row.id) {
        return h('div', { class: 'flex gap-2' }, [
          h(NButton, { size: 'small', type: 'primary', quaternary: true, onClick: () => saveEdit(row) }, { default: () => '保存' }),
          h(NButton, { size: 'small', quaternary: true, onClick: () => { editingId.value = null } }, { default: () => '取消' })
        ])
      }
      return h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => startEdit(row) }, { default: () => '编辑' })
    }
  }
]
function startEdit(row: any) { editingId.value = row.id; editingValue.value = row.value }
async function saveEdit(row: any) {
  try { await put(`/admin/configs/${row.id}`, { value: editingValue.value }); message.success('保存成功'); editingId.value = null; fetchItems() } catch {}
}
async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>('/admin/configs')
    items.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
  } catch { items.value = [] } finally { loading.value = false }
}
onMounted(fetchItems)
</script>