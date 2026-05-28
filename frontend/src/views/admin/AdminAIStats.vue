<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">AI 使用统计</h1>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-400 mb-1">总请求数</p>
        <p class="text-2xl font-bold text-dark">{{ stats.total_requests }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-400 mb-1">总 Token 数</p>
        <p class="text-2xl font-bold text-dark">{{ stats.total_tokens }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-400 mb-1">总错误数</p>
        <p class="text-2xl font-bold text-red-500">{{ stats.total_errors }}</p>
      </div>
    </div>

    <div class="bg-white rounded-xl p-5 border border-gray-100 mb-6">
      <h3 class="text-sm font-semibold text-dark mb-4">近30天使用趋势</h3>
      <div ref="chartRef" style="height: 300px;"></div>
    </div>

    <n-data-table :columns="columns" :data="stats.daily || []" :loading="loading" :bordered="false" :single-line="false" :pagination="{ pageSize: 10 }" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { get } from '@/api/request'
import * as echarts from 'echarts'

const loading = ref(false)
const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null

const stats = ref<any>({
  total_requests: 0,
  total_tokens: 0,
  total_errors: 0,
  daily: []
})

const columns = [
  { title: '日期', key: 'date', width: 120 },
  { title: '请求数', key: 'total_requests', width: 100 },
  { title: 'Token 数', key: 'total_tokens', width: 120 },
  { title: '错误数', key: 'error_count', width: 100 }
]

async function fetchStats() {
  loading.value = true
  try {
    const res = await get<any>('/admin/ai-stats')
    stats.value = res.data
    await nextTick()
    renderChart()
  } catch {} finally { loading.value = false }
}

function renderChart() {
  if (!chartRef.value) return
  if (!chart) {
    chart = echarts.init(chartRef.value)
  }
  const daily = (stats.value.daily || []).slice().reverse()
  chart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: daily.map((d: any) => d.date) },
    yAxis: [
      { type: 'value', name: '请求数' },
      { type: 'value', name: 'Token 数', position: 'right' }
    ],
    series: [
      { name: '请求数', type: 'bar', data: daily.map((d: any) => d.total_requests), itemStyle: { color: '#EB9463', borderRadius: [4, 4, 0, 0] } },
      { name: 'Token 数', type: 'line', yAxisIndex: 1, data: daily.map((d: any) => d.total_tokens), smooth: true, lineStyle: { color: '#6366f1' }, itemStyle: { color: '#6366f1' } }
    ]
  })
}

function handleResize() { chart?.resize() }

onMounted(() => {
  fetchStats()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chart?.dispose()
})
</script>