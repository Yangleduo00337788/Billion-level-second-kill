<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">敏感词管理</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <n-button type="primary" @click="showAdd = true">添加敏感词</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />
    <n-modal v-model:show="showAdd" preset="card" title="添加敏感词" style="max-width: 400px">
      <n-form label-placement="left" label-width="80">
        <n-form-item label="敏感词"><n-input v-model:value="form.word" placeholder="请输入敏感词" /></n-form-item>
        <n-form-item label="级别"><n-select v-model:value="form.level" :options="[{label:'禁止',value:1},{label:'替换',value:2}]" /></n-form-item>
        <n-form-item v-if="form.level===2" label="替换为"><n-input v-model:value="form.replace_to" placeholder="替换内容" /></n-form-item>
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
const form = ref({ word: '', level: 1, replace_to: '' })
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '敏感词', key: 'word', width: 150 },
  { title: '级别', key: 'level', width: 80, render: (row: any) => h(NTag, { type: row.level === 1 ? 'error' : 'warning', size: 'small' }, { default: () => row.level === 1 ? '禁止' : '替换' }) },
  { title: '替换为', key: 'replace_to', width: 120 },
  { title: '时间', key: 'created_at', width: 180 },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row.id) }, { default: () => '删除' }) }
]
async function fetchItems() {
  loading.value = true
  try { const res = await get<any>('/admin/sensitive-words'); items.value = res.data || [] } catch {} finally { loading.value = false }
}
async function handleAdd() {
  if (!form.value.word) { message.warning('请输入敏感词'); return }
  try { await post('/admin/sensitive-words', form.value); message.success('添加成功'); showAdd.value = false; form.value = { word: '', level: 1, replace_to: '' }; fetchItems() } catch {}
}
async function handleDelete(id: number) {
  try { await del(`/admin/sensitive-words/${id}`); message.success('删除成功'); fetchItems() } catch {}
}
onMounted(fetchItems)
</script>