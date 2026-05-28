<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">访问统计</h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-xl p-5 border border-gray-100 hover:shadow-sm transition-shadow">
        <p class="text-sm text-gray-400 mb-1">总 PV</p>
        <p class="text-2xl font-bold text-dark">{{ stats.total_pv || 0 }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100 hover:shadow-sm transition-shadow">
        <p class="text-sm text-gray-400 mb-1">今日 PV</p>
        <p class="text-2xl font-bold text-primary">{{ stats.today_pv || 0 }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100 hover:shadow-sm transition-shadow">
        <p class="text-sm text-gray-400 mb-1">总 UV（独立IP）</p>
        <p class="text-2xl font-bold text-dark">{{ stats.total_uv || 0 }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100 hover:shadow-sm transition-shadow">
        <p class="text-sm text-gray-400 mb-1">今日 UV</p>
        <p class="text-2xl font-bold text-primary">{{ stats.today_uv || 0 }}</p>
      </div>
    </div>
    <div class="bg-white rounded-xl border border-gray-100 p-5">
      <h3 class="text-sm font-semibold text-dark mb-4">热门页面 TOP 10</h3>
      <div v-if="loading" class="text-center text-gray-400 py-8">加载中...</div>
      <div v-else-if="(stats.top_pages || []).length === 0" class="text-center text-gray-400 py-8">暂无数据，页面访问会自动记录</div>
      <n-data-table v-else :columns="columns" :data="stats.top_pages || []" :bordered="false" :single-line="false" />
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { get } from '@/api/request'

const stats = ref<any>({ total_pv: 0, today_pv: 0, total_uv: 0, today_uv: 0, top_pages: [] })
const loading = ref(true)

const columns = [
  {
    title: '页面路径', key: 'path', minWidth: 200,
    render: (row: any) => h('span', { class: 'text-sm font-mono' }, row.path)
  },
  {
    title: '访问量', key: 'count', width: 120,
    render: (row: any) => h('span', { class: 'text-sm font-bold text-primary' }, String(row.count))
  }
]

async function fetchStats() {
  loading.value = true
  try {
    const res = await get<any>('/admin/page-view-stats')
    stats.value = res.data || {}
  } catch {} finally { loading.value = false }
}
onMounted(fetchStats)
</script>