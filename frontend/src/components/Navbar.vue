<template>
  <nav class="fixed top-0 left-0 right-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-100">
    <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
      <div class="flex items-center gap-8">
        <button class="md:hidden p-2 -ml-2 text-gray-500 hover:text-dark transition-colors" @click="mobileMenuOpen = !mobileMenuOpen">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <router-link to="/" class="text-lg font-semibold text-dark tracking-tight">
          推理引擎
        </router-link>
        <div class="hidden md:flex items-center gap-6">
          <router-link
            v-for="link in navLinks"
            :key="link.path"
            :to="link.path"
            class="text-sm text-gray-500 hover:text-dark transition-colors"
            :class="{ 'text-dark font-medium': $route.path === link.path }"
          >
            {{ link.name }}
          </router-link>
        </div>
      </div>

      <div class="flex items-center gap-4">
        <button class="sm:hidden p-2 text-gray-500 hover:text-dark transition-colors" @click="mobileSearchOpen = !mobileSearchOpen">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>
        <div class="hidden sm:flex items-center bg-gray-50 rounded-lg px-3 py-1.5 border border-gray-100">
          <svg class="w-4 h-4 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索文章..."
            class="bg-transparent text-sm outline-none w-40 text-gray-600 placeholder-gray-400"
            @keyup.enter="doSearch"
          />
        </div>

        <template v-if="userStore.isAuthenticated && userStore.user">
          <n-dropdown trigger="click" :options="menuOptions" @select="handleMenuSelect">
            <n-avatar
              :src="userStore.user.avatar"
              round
              size="small"
              class="cursor-pointer hover:opacity-80 transition-opacity"
            />
          </n-dropdown>
        </template>
        <template v-else>
          <router-link to="/login">
            <n-button quaternary size="small" class="text-sm">登录</n-button>
          </router-link>
          <router-link to="/register">
            <n-button size="small" type="primary" class="text-sm">注册</n-button>
          </router-link>
        </template>
      </div>
    </div>

    <div v-if="mobileSearchOpen" class="sm:hidden px-4 py-3 border-b border-gray-100 bg-white">
      <div class="flex items-center bg-gray-50 rounded-lg px-3 py-2 border border-gray-100">
        <svg class="w-4 h-4 text-gray-400 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
        <input v-model="searchQuery" type="text" placeholder="搜索..." class="bg-transparent text-sm outline-none w-full text-gray-600 placeholder-gray-400" @keyup.enter="doSearch; mobileSearchOpen = false" />
      </div>
    </div>

    <div v-if="mobileMenuOpen" class="md:hidden absolute top-16 left-0 right-0 bg-white border-b border-gray-100 shadow-lg">
      <div class="px-4 py-3">
        <div class="flex flex-col gap-1">
          <router-link v-for="link in navLinks" :key="link.path" :to="link.path" class="px-3 py-2 text-sm text-gray-600 hover:text-dark hover:bg-gray-50 rounded-lg transition-colors" @click="mobileMenuOpen = false">{{ link.name }}</router-link>
        </div>
        <div class="mt-2 pt-2 border-t border-gray-100">
          <template v-if="userStore.isAuthenticated">
            <router-link to="/editor" class="block px-3 py-2 text-sm text-gray-600 hover:text-dark hover:bg-gray-50 rounded-lg" @click="mobileMenuOpen = false">写文章</router-link>
            <button class="block w-full text-left px-3 py-2 text-sm text-gray-600 hover:text-dark hover:bg-gray-50 rounded-lg" @click="handleLogout">退出</button>
          </template>
          <template v-else>
            <router-link to="/login" class="block px-3 py-2 text-sm text-primary font-medium hover:bg-gray-50 rounded-lg" @click="mobileMenuOpen = false">登录</router-link>
            <router-link to="/register" class="block px-3 py-2 text-sm text-primary font-medium hover:bg-gray-50 rounded-lg" @click="mobileMenuOpen = false">注册</router-link>
          </template>
        </div>
      </div>
    </div>
  </nav>
  <div class="h-16"></div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NDropdown, NAvatar } from 'naive-ui'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

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
</script>
