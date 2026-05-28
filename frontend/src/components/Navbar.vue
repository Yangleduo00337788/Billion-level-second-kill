<template>
  <nav class="fixed top-0 left-0 right-0 z-50 nav-glass">
    <div class="max-w-[1440px] mx-auto px-6 lg:px-10 h-16 flex items-center justify-between">
      <div class="flex items-center gap-8">
        <button class="md:hidden p-2 -ml-2 text-gray-500 hover:text-dark transition-colors rounded-xl hover:bg-white/30" @click="mobileMenuOpen = !mobileMenuOpen">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <router-link to="/" class="glass-logo">
          <svg class="w-8 h-8 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456z" />
          </svg>
          <span class="text-lg font-bold text-dark tracking-tight relative z-10">{{ siteName }}</span>
        </router-link>
        <div class="hidden md:flex items-center gap-1">
          <router-link
            v-for="link in navLinks"
            :key="link.path"
            :to="link.path"
            class="glass-nav-link text-sm"
            :class="$route.path === link.path
              ? 'active text-primary font-medium'
              : 'text-gray-600 hover:text-dark'"
          >
            {{ link.name }}
          </router-link>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button class="sm:hidden p-2 text-gray-500 hover:text-dark transition-colors rounded-xl hover:bg-white/30" @click="mobileSearchOpen = !mobileSearchOpen">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>
        <div class="hidden sm:flex items-center glass-input px-3 py-2">
          <svg class="w-4 h-4 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索文章..."
            class="bg-transparent text-sm outline-none w-40 text-gray-700 placeholder-gray-400"
            @keyup.enter="doSearch"
          />
        </div>

        <!-- Announcement Icon -->
        <n-popover trigger="click" :show="showAnnouncements" @update:show="showAnnouncements = $event" placement="bottom-end" :width="320">
          <template #trigger>
            <button class="relative p-2 text-gray-500 hover:text-dark transition-colors rounded-xl hover:bg-white/30">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
              </svg>
              <span v-if="unreadCount > 0" class="absolute -top-0.5 -right-0.5 w-4 h-4 bg-red-500 text-white text-[10px] rounded-full flex items-center justify-center font-medium">{{ unreadCount > 9 ? '9+' : unreadCount }}</span>
            </button>
          </template>
          <div class="max-h-80 overflow-y-auto">
            <div class="px-3 py-2 border-b border-gray-100 flex items-center justify-between">
              <span class="text-sm font-semibold text-dark">公告通知</span>
              <span class="text-xs text-gray-400">共 {{ announcements.length }} 条</span>
            </div>
            <div v-if="announcements.length === 0" class="py-8 text-center text-gray-400 text-sm">暂无公告</div>
            <div v-else>
              <div
                v-for="item in announcements"
                :key="item.id"
                class="px-3 py-3 border-b border-gray-50 hover:bg-gray-50 transition-colors cursor-pointer"
                @click="openAnnouncement(item)"
              >
                <div class="flex items-start gap-2">
                  <span class="flex-shrink-0 w-2 h-2 rounded-full mt-1.5" :class="getPriorityColor(item.priority)"></span>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-dark truncate">{{ item.title }}</p>
                    <p class="text-xs text-gray-500 mt-1 line-clamp-2">{{ item.content }}</p>
                    <p class="text-xs text-gray-400 mt-1">{{ item.created_at }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </n-popover>

        <!-- Theme Toggle -->
        <n-dropdown trigger="click" :options="themeOptions" @select="handleThemeSelect">
          <button class="p-2 text-gray-500 hover:text-dark transition-colors rounded-xl hover:bg-white/30 dark:hover:bg-white/10" :title="themeLabel">
            <svg v-if="themeStore.mode === 'light'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else-if="themeStore.mode === 'dark'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </button>
        </n-dropdown>

        <template v-if="userStore.isAuthenticated && userStore.user">
          <div class="flex items-center gap-1 px-2 py-1 bg-yellow-50 rounded-full" :title="`${userStore.user.points || 0} 积分`">
            <svg class="w-4 h-4 text-yellow-500" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
            </svg>
            <span class="text-sm font-medium text-yellow-700">{{ userStore.user.points || 0 }}</span>
          </div>
          <router-link v-if="userStore.user?.role === 'admin'" to="/admin" class="glass-nav-link text-sm text-gray-600 hover:text-dark">
            ⚙ 管理
          </router-link>
          <n-dropdown trigger="click" :options="menuOptions" @select="handleMenuSelect">
            <n-avatar
              :src="userStore.user.avatar"
              round
              size="small"
              class="cursor-pointer hover:opacity-80 transition-opacity ring-2 ring-white/50 hover:ring-primary/30"
            />
          </n-dropdown>
        </template>
        <template v-else>
          <router-link to="/login">
            <button class="glass-button px-4 py-2 text-sm text-gray-700">
              登录
            </button>
          </router-link>
          <router-link to="/register">
            <button class="glass-button-primary px-4 py-2 text-sm">
              注册
            </button>
          </router-link>
        </template>
      </div>
    </div>

    <div v-if="mobileSearchOpen" class="sm:hidden px-4 py-3 border-t border-white/30">
      <div class="flex items-center glass-input px-3 py-2">
        <svg class="w-4 h-4 text-gray-400 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索..."
          class="bg-transparent text-sm outline-none w-full text-gray-700 placeholder-gray-400"
          @keyup.enter="doSearch; mobileSearchOpen = false"
        />
      </div>
    </div>

    <div v-if="mobileMenuOpen" class="md:hidden absolute top-16 left-0 right-0 glass-card border-t border-white/30 shadow-lg rounded-none">
      <div class="px-4 py-3">
        <div class="flex flex-col gap-1">
          <router-link
            v-for="link in navLinks"
            :key="link.path"
            :to="link.path"
            class="px-4 py-2.5 text-sm rounded-xl transition-all duration-200"
            :class="$route.path === link.path
              ? 'text-primary font-medium bg-primary/10'
              : 'text-gray-700 hover:text-dark hover:bg-white/30'"
            @click="mobileMenuOpen = false"
          >
            {{ link.name }}
          </router-link>
        </div>
        <div class="mt-3 pt-3 border-t border-white/30">
          <template v-if="userStore.isAuthenticated">
            <router-link to="/editor" class="block px-4 py-2.5 text-sm text-gray-700 hover:text-dark hover:bg-white/30 rounded-xl" @click="mobileMenuOpen = false">
              写文章
            </router-link>
            <button class="block w-full text-left px-4 py-2.5 text-sm text-gray-700 hover:text-dark hover:bg-white/30 rounded-xl" @click="handleLogout">
              退出
            </button>
          </template>
          <template v-else>
            <router-link to="/login" class="block px-4 py-2.5 text-sm text-primary font-medium hover:bg-primary/5 rounded-xl" @click="mobileMenuOpen = false">
              登录
            </router-link>
            <router-link to="/register" class="block px-4 py-2.5 text-sm text-primary font-medium hover:bg-primary/5 rounded-xl" @click="mobileMenuOpen = false">
              注册
            </router-link>
          </template>
        </div>
      </div>
    </div>
  </nav>
  <div class="h-16"></div>

  <!-- Announcement Detail Modal -->
  <n-modal v-model:show="showDetailModal" preset="card" :title="currentAnnouncement?.title" style="max-width: 500px">
    <p class="text-sm text-gray-600 whitespace-pre-wrap">{{ currentAnnouncement?.content }}</p>
    <template #footer>
      <div class="text-xs text-gray-400">{{ currentAnnouncement?.created_at }}</div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NDropdown, NAvatar, NPopover, NModal } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { get } from '@/api/request'
import type { ThemeMode } from '@/stores/theme'

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()

const announcements = ref<any[]>([])
const showAnnouncements = ref(false)
const showDetailModal = ref(false)
const currentAnnouncement = ref<any>(null)
const unreadCount = ref(0)
const siteName = ref('推理引擎')

async function fetchSiteConfig() {
  try {
    const res = await get<any>('/site-config')
    if (res.data?.site_name) {
      siteName.value = res.data.site_name
    }
  } catch {}
}

const themeOptions = [
  { key: 'light', label: '浅色模式' },
  { key: 'dark', label: '深色模式' },
  { key: 'system', label: '跟随系统' }
]

const themeLabel = computed(() => {
  const map = { light: '浅色模式', dark: '深色模式', system: '跟随系统' }
  return map[themeStore.mode]
})

function handleThemeSelect(key: string) {
  themeStore.setMode(key as ThemeMode)
}

const searchQuery = ref('')
const mobileMenuOpen = ref(false)
const mobileSearchOpen = ref(false)

const navLinks = [
  { name: '首页', path: '/' },
  { name: '发现', path: '/search' },
  { name: 'Prompt', path: '/prompt' },
  { name: 'AI 助手', path: '/chat' },
  { name: '消息', path: '/im' }
]

const menuOptions = [
  { key: 'profile', label: '我的主页' },
  { key: 'editor', label: '写文章' },
  { key: 'prompts', label: '我的 Prompt' },
  { type: 'divider' as const },
  { key: 'logout', label: '退出' }
]

function handleMenuSelect(key: string) {
  switch (key) {
    case 'profile':
      router.push(`/user/${userStore.user?.id}`)
      break
    case 'editor':
      router.push('/editor')
      break
    case 'prompts':
      router.push(`/user/${userStore.user?.id}?tab=prompts`)
      break
    case 'logout':
      userStore.logout()
      router.push('/')
      break
  }
}

function doSearch() {
  if (searchQuery.value.trim()) {
    router.push({ path: '/search', query: { q: searchQuery.value.trim() } })
  }
}

function handleLogout() {
  userStore.logout()
  mobileMenuOpen.value = false
  router.push('/')
}

function getPriorityColor(priority: number) {
  return ['', 'bg-yellow-500', 'bg-red-500'][priority] || 'bg-gray-300'
}

function openAnnouncement(item: any) {
  currentAnnouncement.value = item
  showDetailModal.value = true
  showAnnouncements.value = false
}

async function fetchAnnouncements() {
  try {
    const res = await get<any>('/announcements')
    const items = Array.isArray(res.data) ? res.data : (res.data?.items || [])
    announcements.value = items
    unreadCount.value = items.filter((i: any) => i.priority >= 1).length
  } catch {}
}

onMounted(() => {
  fetchAnnouncements()
  fetchSiteConfig()
})
</script>