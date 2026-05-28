<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <header class="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-6 flex-shrink-0">
      <div class="flex items-center gap-4">
        <router-link to="/admin" class="flex items-center gap-2">
          <svg class="w-8 h-8 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456z" />
          </svg>
          <span class="text-lg font-bold text-dark">推理引擎</span>
          <span class="text-xs bg-primary/10 text-primary px-2 py-0.5 rounded-full font-medium">Admin</span>
        </router-link>
      </div>
      <div class="flex items-center gap-4">
        <router-link to="/" class="text-sm text-gray-500 hover:text-dark transition-colors flex items-center gap-1.5">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
          返回前台
        </router-link>
        <div class="w-px h-5 bg-gray-200"></div>
        <n-avatar :src="userStore.user?.avatar" round size="small" />
        <span class="text-sm text-gray-600">{{ userStore.user?.username }}</span>
      </div>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <aside class="w-60 bg-white border-r border-gray-200 flex-shrink-0 hidden md:flex flex-col overflow-y-auto">
        <nav class="flex-1 p-3 space-y-0.5">
          <template v-for="item in menuItems" :key="item.path || item.key">
            <div v-if="item.children" class="mb-1">
              <button @click="toggleGroup(item.key)" class="w-full flex items-center justify-between px-3 py-2.5 text-sm text-gray-600 hover:text-dark rounded-xl transition-colors">
                <div class="flex items-center gap-3">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" v-html="item.icon"></svg>
                  <span>{{ item.label }}</span>
                </div>
                <svg class="w-4 h-4 transition-transform" :class="{ '-rotate-90': !expandedGroups[item.key] }" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" /></svg>
              </button>
              <div v-show="expandedGroups[item.key]" class="ml-4 space-y-0.5 mt-1">
                <router-link v-for="child in item.children" :key="child.path" :to="child.path" class="flex items-center gap-3 px-3 py-2 text-sm rounded-xl transition-colors" :class="isActive(child.path) ? 'bg-primary/10 text-primary font-medium' : 'text-gray-500 hover:bg-gray-50 hover:text-dark'">
                  <span class="w-1.5 h-1.5 rounded-full" :class="isActive(child.path) ? 'bg-primary' : 'bg-gray-300'"></span>
                  {{ child.label }}
                </router-link>
              </div>
            </div>
            <router-link v-else :to="item.path" class="flex items-center gap-3 px-3 py-2.5 text-sm rounded-xl transition-colors" :class="isActive(item.path) ? 'bg-primary/10 text-primary font-medium' : 'text-gray-600 hover:bg-gray-50 hover:text-dark'">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" v-html="item.icon"></svg>
              {{ item.label }}
            </router-link>
          </template>
        </nav>
      </aside>

      <main class="flex-1 overflow-y-auto">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRoute } from 'vue-router'
import { NAvatar } from 'naive-ui'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()

const expandedGroups = reactive<Record<string, boolean>>({
  content: false,
  users: false,
  security: false,
  operations: false,
  system: false
})

function toggleGroup(key: string) {
  expandedGroups[key] = !expandedGroups[key]
}

const menuItems = [
  {
    path: '/admin',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />',
    label: '仪表盘'
  },
  {
    key: 'users',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />',
    label: '用户管理',
    children: [
      { path: '/admin/users', label: '用户列表' },
      { path: '/admin/user-tags', label: '用户标签' }
    ]
  },
  {
    key: 'content',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />',
    label: '内容管理',
    children: [
      { path: '/admin/articles', label: '文章管理' },
      { path: '/admin/prompts', label: 'Prompt 管理' },
      { path: '/admin/comments', label: '评论管理' },
      { path: '/admin/categories', label: '分类管理' },
      { path: '/admin/content-reviews', label: '内容审核' },
      { path: '/admin/recommend', label: '推荐位管理' }
    ]
  },
  {
    path: '/admin/announcements',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />',
    label: '公告管理'
  },
  {
    path: '/admin/reports',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />',
    label: '举报管理'
  },
  {
    path: '/admin/rankings',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />',
    label: '排行榜'
  },
  {
    key: 'security',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />',
    label: '安全管理',
    children: [
      { path: '/admin/audit-logs', label: '操作日志' },
      { path: '/admin/login-logs', label: '登录日志' },
      { path: '/admin/sensitive-words', label: '敏感词管理' },
      { path: '/admin/ip-blacklist', label: 'IP 黑名单' }
    ]
  },
  {
    key: 'operations',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />',
    label: '运营管理',
    children: [
      { path: '/admin/points-rules', label: '积分规则' },
      { path: '/admin/invite-codes', label: '邀请码管理' }
    ]
  },
  {
    key: 'system',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />',
    label: '系统管理',
    children: [
      { path: '/admin/ai-stats', label: 'AI 统计' },
      { path: '/admin/page-stats', label: '访问统计' },
      { path: '/admin/system-logs', label: '系统日志' },
      { path: '/admin/configs', label: '系统配置' },
      { path: '/admin/notifications', label: '通知管理' },
      { path: '/admin/settings', label: '系统设置' }
    ]
  }
]

function isActive(path: string) {
  if (path === '/admin') return route.path === '/admin'
  return route.path.startsWith(path)
}

// Auto expand active group based on current route
const currentPath = route.path
menuItems.forEach(item => {
  if (item.children) {
    const hasActive = item.children.some((child: any) => currentPath.startsWith(child.path))
    if (hasActive) expandedGroups[item.key] = true
  }
})
</script>