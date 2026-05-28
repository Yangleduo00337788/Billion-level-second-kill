<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-dark">分类管理</h1>
      <n-button type="primary" @click="showAdd = true">新增分类</n-button>
    </div>
    <div class="bg-white rounded-2xl border border-gray-100 overflow-hidden">
      <n-data-table :columns="columns" :data="categories" :loading="loading" :bordered="false" :row-key="(row: any) => row.id" />
    </div>
    <n-modal v-model:show="showAdd" preset="card" title="新增分类" style="width: 400px;">
      <n-form>
        <n-form-item label="名称"><n-input v-model:value="newCat.name" placeholder="分类名称" /></n-form-item>
        <n-form-item label="描述"><n-input v-model:value="newCat.desc" placeholder="分类描述" /></n-form-item>
      </n-form>
      <template #footer><n-button type="primary" @click="addCategory">确认</n-button></template>
    </n-modal>
  </div>
</template>
<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { get, post, del } from '@/api/request'
import { useMessage, NButton, type DataTableColumns } from 'naive-ui'
const message = useMessage()
const loading = ref(false)
const categories = ref<any[]>([])
const showAdd = ref(false)
const newCat = ref({ name: '', desc: '' })
const columns: DataTableColumns<any> = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '描述', key: 'desc' },
  { title: '排序', key: 'sort', width: 60 },
  { title: '操作', key: 'actions', width: 80, render: (row) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row) }, () => '删除') }
]
async function fetchCategories() {
  loading.value = true
  try { const res = await get('/articles/categories'); categories.value = res.data || [] } catch { categories.value = [] } finally { loading.value = false }
}
async function addCategory() {
  if (!newCat.value.name.trim()) { message.warning('请输入名称'); return }
  try { await post('/admin/categories', newCat.value); message.success('已添加'); showAdd.value = false; newCat.value = { name: '', desc: '' }; fetchCategories() } catch { message.error('添加失败') }
}
function handleDelete(row: any) {
  del('/admin/categories/' + row.id).then(() => { categories.value = categories.value.filter((c: any) => c.id !== row.id); message.success('已删除') }).catch(() => message.error('删除失败'))
}
onMounted(fetchCategories)
</script>