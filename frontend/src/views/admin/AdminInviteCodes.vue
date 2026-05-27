<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">邀请码管理</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <n-button type="primary" @click="showAdd = true">生成邀请码</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />
    <n-modal v-model:show="showAdd" preset="card" title="生成邀请码" style="max-width: 400px">
      <n-form label-placement="left" label-width="100">
        <n-form-item label="最大使用次数"><n-input-number v-model:value="form.max_uses" :min="1" style="width: 100%" /></n-form-item>
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
import { NButton, NTag, NInput, useMessage } from 'naive-ui'
import { get, post } from '@/api/request'
const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])
const showAdd = ref(false)
const form = ref({ max_uses: 10 })
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '邀请码', key: 'code', width: 150, render: (row: any) => h('span', { class: 'font-mono text-sm bg-gray-100 px-2 py-1 rounded' }, row.code) },
  { title: '最大使用', key: 'max_uses', width: 100 },
  { title: '已使用', key: 'used_count', width: 80 },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: (row.used_count || 0) < (row.max_uses || 1) ? 'success' : 'default', size: 'small' }, { default: () => (row.used_count || 0) < (row.max_uses || 1) ? '可用' : '已用完' }) },
  { title: '过期时间', key: 'expires_at', width: 170, render: (row: any) => row.expires_at || '永不过期' },
  { title: '创建时间', key: 'created_at', width: 170 }
]
async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>('/admin/invite-codes')
    items.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
  } catch { items.value = [] } finally { loading.value = false }
}
async function handleAdd() {
  try {
    await post('/admin/invite-codes', form.value)
    message.success('生成成功')
    showAdd.value = false
    fetchItems()
  } catch {}
}
onMounted(fetchItems)
</script>