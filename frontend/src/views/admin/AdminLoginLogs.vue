<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-dark mb-6">登录日志</h1>
    <n-data-table :columns="columns" :data="items" :loading="loading" :bordered="false" :single-line="false" :pagination="pagination" remote />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted, reactive, h } from 'vue'
import { NTag } from 'naive-ui'
import { get } from '@/api/request'

const loading = ref(false)
const items = ref<any[]>([])
const pagination = reactive({
  page: 1, pageSize: 20, pageCount: 1,
  showSizePicker: true, pageSizes: [20, 50, 100],
  onChange: (page: number) => { pagination.page = page; fetchItems() },
  onUpdatePageSize: (size: number) => { pagination.pageSize = size; pagination.page = 1; fetchItems() }
})

function DeviceIcon({ type }: { type: string }) {
  if (type === '手机') {
    return h('svg', { class: 'w-4 h-4 inline-block text-gray-500', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', innerHTML: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />' })
  }
  if (type === '平板') {
    return h('svg', { class: 'w-4 h-4 inline-block text-gray-500', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', innerHTML: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />' })
  }
  return h('svg', { class: 'w-4 h-4 inline-block text-gray-500', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', innerHTML: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />' })
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '用户名', key: 'username', width: 100,
    render: (row: any) => h('span', { class: 'text-sm font-medium' }, row.username || (row.user_id ? `#${row.user_id}` : '-'))
  },
  {
    title: '状态', key: 'status', width: 80,
    render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small', round: true }, { default: () => row.status === 1 ? '成功' : '失败' })
  },
  {
    title: 'IP 地址', key: 'ip', width: 140,
    render: (row: any) => {
      const ip = row.ip || '-'
      const display = ip === '::1' ? '127.0.0.1 (本机)' : ip
      return h('span', { class: 'text-xs font-mono text-gray-500' }, display)
    }
  },
  {
    title: '设备', key: 'device_type', width: 80,
    render: (row: any) => h('span', { class: 'flex items-center gap-1 text-sm' }, [
      h(DeviceIcon, { type: row.device_type || '电脑' }),
      h('span', {}, row.device_type || '未知')
    ])
  },
  {
    title: '浏览器', key: 'browser', width: 80,
    render: (row: any) => h('span', { class: 'text-sm' }, row.browser || '未知')
  },
  {
    title: '操作系统', key: 'os', width: 90,
    render: (row: any) => h('span', { class: 'text-sm' }, row.os || '未知')
  },
  {
    title: '登录时间', key: 'created_at', width: 170,
    render: (row: any) => h('span', { class: 'text-xs text-gray-500' }, row.created_at ? new Date(row.created_at).toLocaleString('zh-CN') : '-')
  },
  {
    title: '离线时间', key: 'logout_at', width: 170,
    render: (row: any) => {
      if (!row.logout_at) return h(NTag, { type: 'info', size: 'small' }, { default: () => '在线中' })
      return h('span', { class: 'text-xs text-gray-500' }, new Date(row.logout_at).toLocaleString('zh-CN'))
    }
  }
]

async function fetchItems() {
  loading.value = true
  try {
    const res = await get<any>(`/admin/login-logs?page=${pagination.page}&page_size=${pagination.pageSize}`)
    const data = res.data
    items.value = data?.items || data?.data || (Array.isArray(data) ? data : [])
    pagination.pageCount = Math.ceil((data?.total || 0) / pagination.pageSize)
  } catch { items.value = [] } finally { loading.value = false }
}
onMounted(fetchItems)
</script>