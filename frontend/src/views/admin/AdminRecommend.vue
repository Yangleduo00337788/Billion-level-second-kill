<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">推荐位管理</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <n-button type="primary" @click="showAdd = true">添加推荐</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />
    <n-modal v-model:show="showAdd" preset="card" title="添加推荐位" style="max-width: 400px">
      <n-form label-placement="left" label-width="80">
        <n-form-item label="类型"><n-select v-model:value="form.target_type" :options="[{label:'文章',value:'article'},{label:'Prompt',value:'prompt'}]" /></n-form-item>
        <n-form-item label="目标ID"><n-input-number v-model:value="form.target_id" :min="1" /></n-form-item>
        <n-form-item label="位置"><n-input v-model:value="form.position" placeholder="如: homepage_top" /></n-form-item>
        <n-form-item label="排序"><n-input-number v-model:value="form.sort_order" :min="0" /></n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showAdd = false">取消</n-button>
        <n-button type="primary" @click="handleAdd">确定</n-button>
      </template>
    </n-modal>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, NTag, useMessage } from 'naive-ui'
import { get, post, del } from '@/api/request'
const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])
const showAdd = ref(false)
const form = ref({ target_type: 'article', target_id: 1, position: 'homepage_top', sort_order: 0 })
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '类型', key: 'target_type', width: 80, render: (row: any) => h(NTag, { size: 'small' }, { default: () => row.target_type }) },
  { title: '目标ID', key: 'target_id', width: 80 },
  { title: '位置', key: 'position', width: 150 },
  { title: '排序', key: 'sort_order', width: 80 },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small' }, { default: () => row.status === 1 ? '启用' : '禁用' }) },
  { title: '操作', key: 'actions', width: 80, render: (row: any) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row.id) }, { default: () => '删除' }) }
]
async function fetchItems() {
  loading.value = true
  try { const res = await get<any>('/admin/recommend-items'); items.value = Array.isArray(res.data) ? res.data : [] } catch {} finally { loading.value = false }
}
async function handleAdd() {
  try { await post('/admin/recommend-items', form.value); message.success('添加成功'); showAdd.value = false; fetchItems() } catch {}
}
async function handleDelete(id: number) {
  try { await del(`/admin/recommend-items/${id}`); message.success('删除成功'); fetchItems() } catch {}
}
onMounted(fetchItems)
</script>