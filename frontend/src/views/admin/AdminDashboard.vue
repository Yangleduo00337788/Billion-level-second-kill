<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">仪表盘</h1>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div v-for="card in statCards" :key="card.label" class="bg-white rounded-xl p-5 border border-gray-100 hover:shadow-sm transition-shadow">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 rounded-lg flex items-center justify-center" :class="card.bgClass">
            <svg class="w-5 h-5" :class="card.iconClass" fill="none" stroke="currentColor" viewBox="0 0 24 24" v-html="card.icon"></svg>
          </div>
          <p class="text-sm text-gray-400">{{ card.label }}</p>
        </div>
        <p class="text-2xl font-bold text-dark">{{ card.value }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4">近30天趋势</h3>
        <div ref="trendChartRef" style="height: 300px;"></div>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4">今日数据</h3>
        <div class="space-y-4">
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-500">今日注册用户</span>
            <span class="text-lg font-bold text-primary">{{ stats.today_users }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-500">今日发布文章</span>
            <span class="text-lg font-bold text-primary">{{ stats.today_articles }}</span>
          </div>
          <n-divider />
          <h4 class="text-sm font-semibold text-dark">快捷操作</h4>
          <div class="grid grid-cols-2 gap-2">
            <n-button block @click="$router.push('/admin/users')">用户管理</n-button>
            <n-button block @click="$router.push('/admin/articles')">文章管理</n-button>
            <n-button block @click="$router.push('/admin/prompts')">Prompt 管理</n-button>
            <n-button block @click="$router.push('/admin/announcements')">发布公告</n-button>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4">内容分布</h3>
        <div ref="pieChartRef" style="height: 250px;"></div>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <h3 class="text-sm font-semibold text-dark mb-4">最近公告</h3>
        <div class="space-y-3">
          <div v-for="item in announcements.slice(0, 5)" :key="item.id" class="flex items-start gap-3">
            <span class="flex-shrink-0 w-6 h-6 rounded-full bg-primary/10 text-primary flex items-center justify-center">
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" /></svg>
            </span>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-dark truncate">{{ item.title }}</p>
              <p class="text-xs text-gray-400">{{ formatDate(item.created_at) }}</p>
            </div>
          </div>
          <p v-if="!announcements.length" class="text-sm text-gray-400 text-center py-4">暂无公告</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { NDivider } from 'naive-ui'
import { get } from '@/api/request'
import * as echarts from 'echarts'

const trendChartRef = ref<HTMLElement | null>(null)
const pieChartRef = ref<HTMLElement | null>(null)
let trendChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null

const stats = ref<any>({
  total_users: 0, total_articles: 0, total_comments: 0, total_prompts: 0,
  today_users: 0, today_articles: 0
})
const announcements = ref<any[]>([])
const statCards = ref<any[]>([])

function formatDate(date: string) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

async function fetchStats() {
  try {
    const res = await get<any>('/admin/dashboard')
    stats.value = res.data
    statCards.value = [
      { label: '总用户', value: res.data.total_users, bgClass: 'bg-blue-50', iconClass: 'text-blue-500', icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />' },
      { label: '总文章', value: res.data.total_articles, bgClass: 'bg-green-50', iconClass: 'text-green-500', icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />' },
      { label: '总评论', value: res.data.total_comments, bgClass: 'bg-purple-50', iconClass: 'text-purple-500', icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />' },
      { label: '总 Prompt', value: res.data.total_prompts, bgClass: 'bg-orange-50', iconClass: 'text-orange-500', icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />' }
    ]
  } catch {}
}

async function fetchChartData() {
  try {
    const res = await get<any>('/admin/chart-data')
    const data = res.data || []
    await nextTick()

    if (trendChartRef.value) {
      trendChart = echarts.init(trendChartRef.value)
      trendChart.setOption({
        tooltip: { trigger: 'axis' },
        legend: { data: ['用户', '文章', '评论'], bottom: 0 },
        grid: { left: '3%', right: '4%', bottom: '12%', containLabel: true },
        xAxis: {
          type: 'category',
          data: data.map((d: any) => d.date),
          axisLabel: { formatter: (v: string) => v.slice(5) }
        },
        yAxis: { type: 'value' },
        series: [
          { name: '用户', type: 'line', data: data.map((d: any) => d.users), smooth: true, lineStyle: { color: '#3b82f6' }, itemStyle: { color: '#3b82f6' }, areaStyle: { color: 'rgba(59,130,246,0.08)' } },
          { name: '文章', type: 'line', data: data.map((d: any) => d.articles), smooth: true, lineStyle: { color: '#22c55e' }, itemStyle: { color: '#22c55e' }, areaStyle: { color: 'rgba(34,197,94,0.08)' } },
          { name: '评论', type: 'line', data: data.map((d: any) => d.comments), smooth: true, lineStyle: { color: '#a855f7' }, itemStyle: { color: '#a855f7' }, areaStyle: { color: 'rgba(168,85,247,0.08)' } }
        ]
      })
    }

    if (pieChartRef.value) {
      pieChart = echarts.init(pieChartRef.value)
      const pieData = [
        { value: stats.value.total_articles, name: '文章' },
        { value: stats.value.total_prompts, name: 'Prompt' },
        { value: stats.value.total_comments, name: '评论' }
      ]
      const hasData = pieData.some(d => d.value > 0)
      pieChart.setOption({
        tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
        legend: { bottom: '0' },
        series: [{
          type: 'pie', radius: ['40%', '70%'],
          data: hasData ? pieData : [{ value: 1, name: '暂无数据', itemStyle: { color: '#e5e7eb' } }],
          label: { show: true, formatter: '{b}: {c}' },
          emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.2)' } }
        }]
      })
    }
  } catch {}
}

async function fetchAnnouncements() {
  try {
    const res = await get<any>('/admin/announcements')
    announcements.value = res.data?.items || res.data || []
  } catch {}
}

function handleResize() {
  trendChart?.resize()
  pieChart?.resize()
}

onMounted(() => {
  fetchStats().then(() => {
    fetchChartData()
  })
  fetchAnnouncements()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  pieChart?.dispose()
})
</script>