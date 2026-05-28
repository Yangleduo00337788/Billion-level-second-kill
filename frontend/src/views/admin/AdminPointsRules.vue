<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">积分规则</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4 flex gap-3">
      <n-button v-if="items.length === 0" type="primary" @click="initRules">初始化默认规则</n-button>
      <n-button type="primary" @click="openAdd">自定义添加</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />

    <n-modal v-model:show="showAdd" preset="card" title="添加积分规则" style="max-width: 500px">
      <n-form label-placement="left" label-width="80">
        <n-form-item label="行为类型">
          <n-select v-model:value="addForm.action" :options="actionOptions" filterable tag placeholder="选择或输入行为" />
        </n-form-item>
        <n-form-item label="奖励积分">
          <n-input-number v-model:value="addForm.points" :min="0" :max="9999" class="w-full" />
        </n-form-item>
        <n-form-item label="限制类型">
          <n-select v-model:value="addForm.limit_type" :options="limitTypeOptions" />
        </n-form-item>
        <n-form-item label="说明">
          <n-input v-model:value="addForm.desc" placeholder="描述这个积分规则" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showAdd = false">取消</n-button>
        <n-button type="primary" @click="handleAdd">确定</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showEdit" preset="card" title="编辑积分规则" style="max-width: 500px">
      <n-form label-placement="left" label-width="80">
        <n-form-item label="行为类型">
          <n-input :value="actionLabels[editForm.action] || editForm.action" disabled />
        </n-form-item>
        <n-form-item label="奖励积分">
          <n-input-number v-model:value="editForm.points" :min="0" :max="9999" class="w-full" />
        </n-form-item>
        <n-form-item label="限制类型">
          <n-select v-model:value="editForm.limit_type" :options="limitTypeOptions" />
        </n-form-item>
        <n-form-item label="状态">
          <n-switch v-model:value="editForm.status" :checked-value="1" :unchecked-value="0" />
        </n-form-item>
        <n-form-item label="说明">
          <n-input v-model:value="editForm.desc" placeholder="描述这个积分规则" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showEdit = false">取消</n-button>
        <n-button type="primary" @click="handleEdit">确定</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, NTag, NSwitch, useMessage } from 'naive-ui'
import { get, post, put } from '@/api/request'

const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])
const showAdd = ref(false)
const showEdit = ref(false)
const addForm = ref({ action: '', points: 10, desc: '', limit_type: 'unlimited' })
const editForm = ref({ id: 0, action: '', points: 10, desc: '', limit_type: 'unlimited', status: 1 })

const actionOptions = [
  { label: '用户注册', value: 'register' },
  { label: '发布文章', value: 'publish_article' },
  { label: '发布 Prompt', value: 'publish_prompt' },
  { label: '点赞', value: 'like' },
  { label: '评论', value: 'comment' },
  { label: '关注', value: 'follow' },
  { label: '每日登录', value: 'daily_login' },
  { label: '被收藏', value: 'be_favorited' },
  { label: '文章上热门', value: 'article_hot' }
]

const limitTypeOptions = [
  { label: '无限制', value: 'unlimited' },
  { label: '仅首次', value: 'once' },
  { label: '每日一次', value: 'daily' },
  { label: '每周一次', value: 'weekly' },
  { label: '每月一次', value: 'monthly' }
]

const limitTypeLabels: Record<string, string> = {
  unlimited: '无限制',
  once: '仅首次',
  daily: '每日一次',
  weekly: '每周一次',
  monthly: '每月一次'
}

const actionLabels: Record<string, string> = {
  register: '用户注册奖励',
  publish_article: '发布文章奖励',
  publish_prompt: '发布 Prompt 奖励',
  like: '点赞奖励',
  comment: '评论奖励',
  follow: '关注奖励',
  daily_login: '每日登录奖励',
  be_favorited: '被收藏奖励',
  article_hot: '文章上热门'
}

const defaultRules = [
  { action: 'register', points: 100, desc: '新用户注册', limit_type: 'once' },
  { action: 'publish_article', points: 50, desc: '发布一篇原创文章', limit_type: 'unlimited' },
  { action: 'publish_prompt', points: 30, desc: '发布一个 Prompt', limit_type: 'unlimited' },
  { action: 'like', points: 5, desc: '点赞文章/Prompt', limit_type: 'daily' },
  { action: 'comment', points: 10, desc: '发表评论', limit_type: 'unlimited' },
  { action: 'follow', points: 20, desc: '关注其他用户', limit_type: 'unlimited' },
  { action: 'daily_login', points: 5, desc: '每日首次登录', limit_type: 'daily' }
]

function openAdd() {
  addForm.value = { action: '', points: 10, desc: '', limit_type: 'unlimited' }
  showAdd.value = true
}

function startEdit(row: any) {
  editForm.value = {
    id: row.id,
    action: row.action,
    points: row.points,
    desc: row.desc || '',
    limit_type: row.limit_type || 'unlimited',
    status: row.status ?? 1
  }
  showEdit.value = true
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '行为', key: 'action', width: 140,
    render: (row: any) => h('span', { class: 'text-sm font-medium' }, actionLabels[row.action] || row.action)
  },
  {
    title: '奖励积分', key: 'points', width: 100,
    render: (row: any) => h(NTag, { type: 'warning', size: 'small', round: true }, { default: () => `+${row.points}` })
  },
  {
    title: '限制类型', key: 'limit_type', width: 100,
    render: (row: any) => h(NTag, { type: 'info', size: 'small' }, { default: () => limitTypeLabels[row.limit_type] || row.limit_type || '无限制' })
  },
  {
    title: '状态', key: 'status', width: 80,
    render: (row: any) => h(NSwitch, { value: row.status === 1, onUpdateValue: (val: boolean) => toggleStatus(row, val) })
  },
  {
    title: '说明', key: 'desc',
    render: (row: any) => h('span', { class: 'text-sm text-gray-600' }, row.desc || '-')
  },
  {
    title: '操作', key: 'actions', width: 80,
    render: (row: any) => h(NButton, { size: 'small', quaternary: true, type: 'primary', onClick: () => startEdit(row) }, { default: () => '编辑' })
  }
]

async function toggleStatus(row: any, val: boolean) {
  try {
    await put(`/admin/points-rules/${row.id}`, { status: val ? 1 : 0 })
    row.status = val ? 1 : 0
    message.success(val ? '已启用' : '已禁用')
  } catch {}
}

async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>('/admin/points-rules')
    items.value = Array.isArray(res.data) ? res.data : (res.data?.items || [])
  } catch { items.value = [] } finally { loading.value = false }
}

async function initRules() {
  try {
    for (const rule of defaultRules) {
      await post('/admin/points-rules', rule)
    }
    message.success('初始化成功')
    fetchItems()
  } catch {}
}

async function handleAdd() {
  if (!addForm.value.action) { message.warning('请选择行为类型'); return }
  try {
    await post('/admin/points-rules', addForm.value)
    message.success('添加成功')
    showAdd.value = false
    fetchItems()
  } catch {}
}

async function handleEdit() {
  try {
    await put(`/admin/points-rules/${editForm.value.id}`, {
      points: editForm.value.points,
      desc: editForm.value.desc,
      limit_type: editForm.value.limit_type,
      status: editForm.value.status
    })
    message.success('更新成功')
    showEdit.value = false
    fetchItems()
  } catch {}
}

onMounted(fetchItems)
</script>
