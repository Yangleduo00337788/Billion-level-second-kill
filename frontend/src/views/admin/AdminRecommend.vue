<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">推荐位管理</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <n-button type="primary" @click="openAdd">添加推荐</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />

    <n-modal v-model:show="showAdd" preset="card" title="添加推荐位" style="max-width: 500px">
      <n-form label-placement="left" label-width="80">
        <n-form-item label="推荐类型">
          <n-select v-model:value="form.target_type" :options="[
            { label: '推荐文章', value: 'article' },
            { label: '推荐 Prompt', value: 'prompt' }
          ]" @update:value="onTypeChange" />
        </n-form-item>
        <n-form-item label="选择内容">
          <n-select
            v-model:value="form.target_id"
            :options="contentOptions"
            :loading="searchLoading"
            filterable
            remote
            :remote-method="searchContent"
            placeholder="输入关键词搜索文章或 Prompt"
            clearable
          />
        </n-form-item>
        <n-form-item label="展示位置">
          <n-select v-model:value="form.position" :options="[
            { label: '首页顶部轮播', value: 'homepage_top' },
            { label: '首页侧边栏', value: 'homepage_sidebar' },
            { label: '文章列表置顶', value: 'article_top' },
            { label: 'Prompt 推荐', value: 'prompt_recommend' }
          ]" />
        </n-form-item>
        <n-form-item label="排序权重">
          <n-input-number v-model:value="form.sort_order" :min="0" :max="999" />
          <span class="text-xs text-gray-400 ml-2">数字越小越靠前</span>
        </n-form-item>
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
const searchLoading = ref(false)
const items = ref<any[]>([])
const showAdd = ref(false)
const contentOptions = ref<any[]>([])

const positionLabels: Record<string, string> = {
  homepage_top: '首页顶部轮播',
  homepage_sidebar: '首页侧边栏',
  article_top: '文章列表置顶',
  prompt_recommend: 'Prompt 推荐'
}

const form = ref({ target_type: 'article', target_id: null as number | null, position: 'homepage_top', sort_order: 0 })

function openAdd() {
  form.value = { target_type: 'article', target_id: null, position: 'homepage_top', sort_order: 0 }
  contentOptions.value = []
  showAdd.value = true
  loadDefaultContent()
}

function onTypeChange() {
  form.value.target_id = null
  contentOptions.value = []
  loadDefaultContent()
}

async function loadDefaultContent() {
  searchLoading.value = true
  try {
    if (form.value.target_type === 'article') {
      const res = await get<any>('/articles?page=1&page_size=50')
      const articles = res.data?.items || []
      contentOptions.value = articles.map((a: any) => ({ label: `#${a.id} ${a.title}`, value: a.id }))
    } else {
      const res = await get<any>('/prompts?page=1&page_size=50')
      const prompts = res.data?.items || []
      contentOptions.value = prompts.map((p: any) => ({ label: `#${p.id} ${p.title}`, value: p.id }))
    }
  } catch {} finally { searchLoading.value = false }
}

async function searchContent(query: string) {
  if (!query || query.length < 1) { loadDefaultContent(); return }
  searchLoading.value = true
  try {
    if (form.value.target_type === 'article') {
      const res = await get<any>(`/articles?page=1&page_size=20&keyword=${encodeURIComponent(query)}`)
      const articles = res.data?.items || []
      contentOptions.value = articles.map((a: any) => ({ label: `#${a.id} ${a.title}`, value: a.id }))
    } else {
      const res = await get<any>(`/prompts?page=1&page_size=20&keyword=${encodeURIComponent(query)}`)
      const prompts = res.data?.items || []
      contentOptions.value = prompts.map((p: any) => ({ label: `#${p.id} ${p.title}`, value: p.id }))
    }
  } catch {} finally { searchLoading.value = false }
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '推荐类型', key: 'target_type', width: 100,
    render: (row: any) => h(NTag, { type: row.target_type === 'article' ? 'success' : 'warning', size: 'small' }, { default: () => row.target_type === 'article' ? '文章' : 'Prompt' })
  },
  { title: '内容ID', key: 'target_id', width: 80 },
  {
    title: '展示位置', key: 'position', width: 140,
    render: (row: any) => h('span', { class: 'text-sm' }, positionLabels[row.position] || row.position)
  },
  {
    title: '排序权重', key: 'sort_order', width: 100,
    render: (row: any) => h(NTag, { type: 'info', size: 'small' }, { default: () => String(row.sort_order) })
  },
  {
    title: '状态', key: 'status', width: 80,
    render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small' }, { default: () => row.status === 1 ? '启用' : '禁用' })
  },
  {
    title: '操作', key: 'actions', width: 80,
    render: (row: any) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row.id) }, { default: () => '删除' })
  }
]

async function fetchItems() {
  loading.value = true
  try { const res = await get<any>('/admin/recommend-items'); items.value = Array.isArray(res.data) ? res.data : [] } catch {} finally { loading.value = false }
}

async function handleAdd() {
  if (!form.value.target_id) { message.warning('请选择推荐内容'); return }
  try {
    await post('/admin/recommend-items', form.value)
    message.success('添加成功')
    showAdd.value = false
    fetchItems()
  } catch {}
}

async function handleDelete(id: number) {
  try { await del(`/admin/recommend-items/${id}`); message.success('删除成功'); fetchItems() } catch {}
}

onMounted(fetchItems)
</script>