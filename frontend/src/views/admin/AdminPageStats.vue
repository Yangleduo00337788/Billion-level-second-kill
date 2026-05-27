<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">访问统计</h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-400 mb-1">总 PV</p>
        <p class="text-2xl font-bold text-dark">{{ stats.total_pv || 0 }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-400 mb-1">今日 PV</p>
        <p class="text-2xl font-bold text-primary">{{ stats.today_pv || 0 }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-400 mb-1">总 UV</p>
        <p class="text-2xl font-bold text-dark">{{ stats.total_uv || 0 }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-400 mb-1">今日 UV</p>
        <p class="text-2xl font-bold text-primary">{{ stats.today_uv || 0 }}</p>
      </div>
    </div>
    <div class="bg-white rounded-xl border border-gray-100 p-5">
      <h3 class="text-sm font-semibold text-dark mb-4">热门页面</h3>
      <div v-if="(stats.top_pages || []).length === 0" class="text-center text-gray-400 py-8">暂无数据</div>
      <n-data-table v-else :columns="columns" :data="stats.top_pages || []" :bordered="false" :single-line="false" />
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get } from '@/api/request'
const stats = ref<any>({ total_pv: 0, today_pv: 0, total_uv: 0, today_uv: 0, top_pages: [] })
const columns = [
  { title: '页面', key: 'path', ellipsis: { tooltip: true } },
  { title: '访问量', key: 'count', width: 120 }
]
async function fetchStats() {
  try { const res = await get<any>('/admin/page-view-stats'); stats.value = res.data || {} } catch {}
}
onMounted(fetchStats)
</script>