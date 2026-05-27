<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">用户标签</h1>
    <div class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <div class="flex gap-3">
        <n-input v-model:value="searchUsername" placeholder="输入用户名搜索" clearable style="width: 200px" @keyup.enter="searchUser" />
        <n-button type="primary" @click="searchUser">搜索用户</n-button>
      </div>
    </div>

    <div v-if="foundUser" class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <n-avatar :src="foundUser.avatar" round size="small" />
          <span class="font-medium">{{ foundUser.username }}</span>
          <span class="text-sm text-gray-400">ID: {{ foundUser.id }}</span>
        </div>
        <n-button type="primary" size="small" @click="showAddTag = true">添加标签</n-button>
      </div>
    </div>

    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />

    <n-modal v-model:show="showAddTag" preset="card" title="添加标签" style="max-width: 400px">
      <n-form label-placement="left" label-width="60">
        <n-form-item label="标签">
          <n-select v-model:value="newTag" :options="tagOptions" placeholder="选择标签" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showAddTag = false">取消</n-button>
        <n-button type="primary" @click="handleAddTag">确定</n-button>
      </template>
    </n-modal>
  </div>
</template>
<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NTag, NAvatar, useMessage } from 'naive-ui'
import { get, post, del } from '@/api/request'
const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])
const searchUsername = ref('')
const foundUser = ref<any>(null)
const showAddTag = ref(false)
const newTag = ref('')
const tagOptions = [
  { label: 'VIP', value: 'VIP' },
  { label: '优质创作者', value: '优质创作者' },
  { label: '活跃用户', value: '活跃用户' },
  { label: '新手', value: '新手' },
  { label: '认证作者', value: '认证作者' }
]
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '标签', key: 'tag', width: 150, render: (row: any) => h(NTag, { type: 'info', size: 'small' }, { default: () => row.tag }) },
  { title: '创建时间', key: 'created_at', width: 180 },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => handleDelete(row.id) }, { default: () => '删除' }) }
]

async function searchUser() {
  if (!searchUsername.value) { message.warning('请输入用户名'); return }
  loading.value = true
  try {
    const res = await get<any>('/user/list')
    const users = res.data || []
    const user = users.find((u: any) => u.username === searchUsername.value)
    if (user) {
      foundUser.value = user
      fetchTags(user.id)
    } else {
      message.error('未找到用户')
      foundUser.value = null
      items.value = []
    }
  } catch { message.error('搜索失败') } finally { loading.value = false }
}

async function fetchTags(userId: number) {
  try {
    const res = await get<any>(`/admin/users/${userId}/tags`)
    items.value = Array.isArray(res.data) ? res.data : []
  } catch { items.value = [] }
}

async function handleAddTag() {
  if (!newTag.value || !foundUser.value) return
  try {
    await post(`/admin/users/${foundUser.value.id}/tags`, { tag: newTag.value })
    message.success('添加成功')
    showAddTag.value = false
    newTag.value = ''
    fetchTags(foundUser.value.id)
  } catch {}
}

async function handleDelete(id: number) {
  try { await del(`/admin/users/0/tags/${id}`); message.success('删除成功'); if (foundUser.value) fetchTags(foundUser.value.id) } catch {}
}
</script>