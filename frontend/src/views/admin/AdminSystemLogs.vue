<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">系统日志</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <n-select v-model:value="levelFilter" :options="[{label:'全部',value:''},{label:'ERROR',value:'error'},{label:'WARN',value:'warn'},{label:'INFO',value:'info'}]" clearable style="width: 150px" @update:value="fetchItems" />
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" :pagination="pagination" remote />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, reactive, h } from 'vue'
import { NTag } from 'naive-ui'
import { get } from '@/api/request'
const loading = ref(false)
const items = ref<any[]>([])
const levelFilter = ref('')
const pagination = reactive({
  page: 1, pageSize: 20, pageCount: 1,
  onChange: (p: number) => { pagination.page = p; fetchItems() }
})
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '级别', key: 'level', width: 80, render: (row: any) => h(NTag, { type: {error:'error',warn:'warning',info:'info'}[row.level] || 'default', size: 'small' }, { default: () => row.level }) },
  { title: '消息', key: 'message', ellipsis: { tooltip: true } },
  { title: '来源', key: 'source', width: 150 },
  { title: '时间', key: 'created_at', width: 170 }
]
async function fetchItems() {
  loading.value = true
  try {
    let url = `/admin/system-logs?page=${pagination.page}&page_size=${pagination.pageSize}`
    if (levelFilter.value) url += `&level=${levelFilter.value}`
    const res = await get<any>(url)
    const data = res.data
    items.value = data?.items || data?.data || (Array.isArray(data) ? data : [])
    pagination.pageCount = Math.ceil((data?.total || 0) / pagination.pageSize)
  } catch { items.value = [] } finally { loading.value = false }
}
onMounted(fetchItems)
</script>