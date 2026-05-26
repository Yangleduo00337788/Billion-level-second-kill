<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">仪表盘</h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div v-for="card in statCards" :key="card.label" class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-400 mb-1">{{ card.label }}</p>
        <p class="text-2xl font-bold text-dark">{{ card.value }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4">今日数据</h3>
        <div class="space-y-3">
          <div class="flex justify-between text-sm">
            <span class="text-gray-500">今日注册用户</span>
            <span class="font-medium">{{ stats.today_users }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-500">今日发布文章</span>
            <span class="font-medium">{{ stats.today_articles }}</span>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4">快速操作</h3>
        <div class="space-y-2">
          <n-button block quaternary @click="$router.push('/admin/users')">用户管理</n-button>
          <n-button block quaternary @click="$router.push('/admin/articles')">文章管理</n-button>
          <n-button block quaternary @click="$router.push('/admin/prompts')">Prompt 管理</n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/api/request'

const stats = ref({
  total_users: 0, total_articles: 0, total_comments: 0, total_prompts: 0,
  today_users: 0, today_articles: 0
})

const statCards = ref<Array<{ label: string; value: number }>>([])

async function fetchStats() {
  try {
    const res = await get<any>('/admin/dashboard')
    stats.value = res.data
    statCards.value = [
      { label: '总用户', value: res.data.total_users },
      { label: '总文章', value: res.data.total_articles },
      { label: '总评论', value: res.data.total_comments },
      { label: '总 Prompt', value: res.data.total_prompts }
    ]
  } catch {}
}

onMounted(fetchStats)
</script>
