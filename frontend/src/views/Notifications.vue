<template>
  <div class="max-w-[800px] mx-auto px-6 lg:px-10 py-8">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-dark">消息通知</h1>
      <n-button size="small" quaternary @click="markAllAsRead" :disabled="notifications.length === 0">
        全部已读
      </n-button>
    </div>

    <div v-if="loading" class="text-center py-20">
      <n-spin size="large" />
    </div>

    <template v-else>
      <div v-if="notifications.length === 0" class="text-center py-20">
        <div class="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
        </div>
        <p class="text-gray-500">暂无通知</p>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="notify in notifications"
          :key="notify.id"
          class="bg-white rounded-xl border border-gray-100 p-4 cursor-pointer hover:shadow-sm transition-shadow"
          :class="{ 'bg-blue-50 border-blue-100': !notify.is_read }"
          @click="handleClick(notify)"
        >
          <div class="flex items-start gap-3">
            <!-- 图标 -->
            <div
              class="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0"
              :class="getIconClass(notify.type)"
            >
              <span class="text-lg">{{ getIcon(notify.type) }}</span>
            </div>

            <!-- 内容 -->
            <div class="flex-1 min-w-0">
              <p class="text-sm text-dark" :class="{ 'font-medium': !notify.is_read }">
                {{ notify.content }}
              </p>
              <p class="text-xs text-gray-400 mt-1">{{ formatDate(notify.created_at) }}</p>
            </div>

            <!-- 未读标记 -->
            <div v-if="!notify.is_read" class="w-2 h-2 rounded-full bg-primary flex-shrink-0 mt-2"></div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="mt-6 flex justify-center">
        <n-pagination
          v-model:page="page"
          :page-size="pageSize"
          :item-count="total"
          @update:page="fetchNotifications"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, NSpin, NButton, NPagination } from 'naive-ui'
import { get, put } from '@/api/request'

const router = useRouter()
const message = useMessage()

interface Notification {
  id: number
  user_id: number
  actor_id: number
  type: string
  content: string
  target_id: number
  target_url: string
  is_read: boolean
  created_at: string
}

const notifications = ref<Notification[]>([])
const loading = ref(true)
const page = ref(1)
const pageSize = 20
const total = ref(0)

async function fetchNotifications() {
  loading.value = true
  try {
    const res = await get<any>('/notifications', { page: page.value, page_size: pageSize })
    notifications.value = res.data?.items || []
    total.value = res.data?.total || 0
  } catch {
    notifications.value = []
  } finally {
    loading.value = false
  }
}

async function markAllAsRead() {
  try {
    await put('/notifications/read-all')
    notifications.value.forEach(n => n.is_read = true)
    message.success('已全部标记为已读')
  } catch {
    message.error('操作失败')
  }
}

async function handleClick(notify: Notification) {
  // 标记为已读
  if (!notify.is_read) {
    try {
      await put(`/notifications/${notify.id}/read`)
      notify.is_read = true
    } catch {}
  }

  // 根据类型跳转
  switch (notify.type) {
    case 'like':
    case 'comment':
      if (notify.target_id) {
        router.push(`/article/${notify.target_id}`)
      }
      break
    case 'follow':
      if (notify.actor_id) {
        router.push(`/user/${notify.actor_id}`)
      }
      break
    case 'points':
      // 积分通知跳转到个人主页
      router.push(`/user/${notify.user_id}`)
      break
    default:
      break
  }
}

function getIconClass(type: string) {
  const classMap: Record<string, string> = {
    'like': 'bg-red-50 text-red-500',
    'comment': 'bg-blue-50 text-blue-500',
    'follow': 'bg-green-50 text-green-500',
    'points': 'bg-yellow-50 text-yellow-500',
    'system': 'bg-gray-50 text-gray-500'
  }
  return classMap[type] || 'bg-gray-50 text-gray-500'
}

function getIcon(type: string) {
  const iconMap: Record<string, string> = {
    'like': '❤️',
    'comment': '💬',
    'follow': '👤',
    'points': '⭐',
    'system': '📢'
  }
  return iconMap[type] || '📢'
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} 天前`

  return date.toLocaleDateString('zh-CN')
}

onMounted(fetchNotifications)
</script>
