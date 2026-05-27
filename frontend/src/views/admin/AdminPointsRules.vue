<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">积分规则</h1>
    <div v-if="items.length === 0 && !loading" class="bg-white rounded-xl border border-gray-100 p-4 mb-4">
      <n-button type="primary" @click="initRules">初始化默认规则</n-button>
    </div>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NButton, NTag, useMessage } from 'naive-ui'
import { get, post } from '@/api/request'
const message = useMessage()
const loading = ref(false)
const items = ref<any[]>([])
const defaultRules = [
  { action: 'register', points: 100, desc: '用户注册' },
  { action: 'publish_article', points: 50, desc: '发布文章' },
  { action: 'publish_prompt', points: 30, desc: '发布 Prompt' },
  { action: 'like', points: 5, desc: '点赞' },
  { action: 'comment', points: 10, desc: '评论' },
  { action: 'follow', points: 20, desc: '关注' },
  { action: 'daily_login', points: 5, desc: '每日登录' }
]
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '操作', key: 'action', width: 150 },
  { title: '积分', key: 'points', width: 100, render: (row: any) => h(NTag, { type: 'warning', size: 'small' }, { default: () => `+${row.points}` }) },
  { title: '说明', key: 'desc' },
  { title: '时间', key: 'created_at', width: 170 }
]
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
onMounted(fetchItems)
</script>