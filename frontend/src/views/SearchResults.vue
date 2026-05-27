<template>
  <div class="max-w-[1100px] mx-auto px-6 lg:px-10 py-8">
    <button
      class="mb-6 inline-flex items-center gap-2 text-sm text-gray-500 hover:text-dark transition-colors glass-button px-3 py-1.5 rounded-xl"
      @click="router.back()"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      返回
    </button>

    <div class="mb-8">
      <div class="flex gap-3">
        <n-input
          v-model:value="query"
          placeholder="搜索文章、用户、Prompt..."
          size="large"
          clearable
          class="flex-1"
          @keyup.enter="doSearch"
        >
          <template #prefix>
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </template>
        </n-input>
        <n-button type="primary" size="large" @click="doSearch" :loading="loading">搜索</n-button>
      </div>
    </div>

    <div class="flex items-center gap-1 mb-6 border-b border-gray-100">
      <n-button
        v-for="tab in tabs"
        :key="tab.key"
        :type="activeTab === tab.key ? 'primary' : 'default'"
        :bordered="false"
        size="small"
        @click="switchTab(tab.key)"
      >
        {{ tab.label }}
        <span v-if="tabCounts[tab.key] > 0" class="ml-1 text-xs opacity-60">({{ tabCounts[tab.key] }})</span>
      </n-button>
    </div>

    <div v-if="!hasSearched && !query">
      <h3 class="text-lg font-semibold text-dark mb-4">热门文章</h3>
      <div class="space-y-4">
        <div
          v-for="item in discoverArticles"
          :key="item.id"
          class="bg-white rounded-2xl border border-gray-100 p-5 cursor-pointer hover:shadow-md transition-shadow"
          @click="router.push('/article/' + item.id)"
        >
          <h3 class="font-semibold text-dark mb-2">{{ item.title }}</h3>
          <p class="text-sm text-gray-500 line-clamp-2 mb-3">{{ item.summary }}</p>
          <div class="flex items-center gap-3 text-xs text-gray-400">
            <span>{{ item.user?.username }}</span>
            <span>{{ formatCount(item.view_count) }} 阅读</span>
          </div>
        </div>
        <div v-if="discoverArticles.length === 0" class="text-center py-12 text-gray-400 text-sm">暂无内容</div>
      </div>
    </div>

    <div v-else-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="bg-white rounded-2xl p-5 animate-pulse border border-gray-100">
        <div class="bg-gray-100 h-4 rounded w-2/3 mb-3"></div>
        <div class="bg-gray-100 h-3 rounded w-full mb-2"></div>
        <div class="bg-gray-100 h-3 rounded w-3/4"></div>
      </div>
    </div>

    <div v-else-if="displayResults.length === 0" class="text-center py-20 text-gray-400">
      <p>没有找到相关结果</p>
    </div>

    <div v-else class="space-y-4">
      <!-- Articles -->
      <template v-if="activeTab === 'articles'">
        <div
          v-for="item in displayResults"
          :key="item.id"
          class="bg-white rounded-2xl border border-gray-100 p-5 cursor-pointer hover:shadow-md transition-shadow"
          @click="router.push('/article/' + item.id)"
        >
          <h3 class="font-semibold text-dark mb-2" v-html="highlightText(item.title)"></h3>
          <p class="text-sm text-gray-500 line-clamp-2 mb-3" v-html="highlightText(item.summary || '')"></p>
          <div class="flex items-center gap-3 text-xs text-gray-400">
            <span>{{ item.username }}</span>
            <span>{{ formatCount(item.view_count) }} 阅读</span>
          </div>
        </div>
      </template>

      <!-- Users -->
      <template v-if="activeTab === 'users'">
        <div
          v-for="item in displayResults"
          :key="item.id"
          class="bg-white rounded-2xl border border-gray-100 p-5 cursor-pointer hover:shadow-md transition-shadow flex items-center gap-4"
          @click="router.push('/user/' + item.id)"
        >
          <n-avatar :src="item.avatar" round :size="48" />
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-dark mb-1" v-html="highlightText(item.username)"></h3>
            <p class="text-sm text-gray-500 line-clamp-1">{{ item.bio || '这个人很懒，什么都没写' }}</p>
            <div class="flex items-center gap-4 mt-1 text-xs text-gray-400">
              <span>{{ item.article_count }} 文章</span>
              <span>{{ item.fans_count }} 粉丝</span>
            </div>
          </div>
        </div>
      </template>

      <!-- Prompts -->
      <template v-if="activeTab === 'prompts'">
        <div
          v-for="item in displayResults"
          :key="item.id"
          class="bg-white rounded-2xl border border-gray-100 p-5 cursor-pointer hover:shadow-md transition-shadow"
          @click="router.push('/prompt/' + item.id)"
        >
          <div class="flex items-center gap-2 mb-2">
            <n-tag v-if="item.category" size="tiny" :bordered="false">{{ item.category }}</n-tag>
            <span class="text-xs text-gray-400">使用 {{ formatCount(item.usage_count) }} 次</span>
          </div>
          <h3 class="font-semibold text-dark mb-2" v-html="highlightText(item.title)"></h3>
          <p class="text-sm text-gray-500 line-clamp-2 mb-3" v-html="highlightText(item.description || '')"></p>
          <div class="flex items-center gap-2 text-xs text-gray-400">
            <n-avatar :src="item.user?.avatar" round size="tiny" />
            <span>{{ item.user?.username }}</span>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get } from '@/api/request'
import { formatCount } from '@/utils/helpers'

const route = useRoute()
const router = useRouter()

const query = ref('')
const activeTab = ref('articles')
const loading = ref(false)
const hasSearched = ref(false)

// Store ALL results from API
const allResults = ref<{ articles: any[]; users: any[]; prompts: any[] }>({ articles: [], users: [], prompts: [] })

const tabs = [
  { key: 'articles', label: '文章' },
  { key: 'users', label: '用户' },
  { key: 'prompts', label: 'Prompt' }
]

const tabCounts = computed(() => {
  const counts: Record<string, number> = { articles: allResults.value.articles.length, users: allResults.value.users.length, prompts: allResults.value.prompts.length }
  return counts
})

const displayResults = computed(() => {
  const key = activeTab.value as 'articles' | 'users' | 'prompts'
  return allResults.value[key] || []
})

function escapeHtml(text: string): string {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

function highlightText(text: string): string {
  if (!text || !query.value) return text || ''
  const safe = escapeHtml(text)
  const q = query.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return safe.replace(new RegExp('(' + q + ')', 'gi'), '<mark class="bg-yellow-100 text-yellow-900 px-0.5 rounded">$1</mark>')
}

function switchTab(key: string) {
  activeTab.value = key
}

async function doSearch() {
  if (!query.value.trim()) return
  loading.value = true
  hasSearched.value = true
  try {
    const res = await get('/search', { q: query.value, page: 1, page_size: 20 })
    const data = res.data as any
    allResults.value = {
      articles: (data.articles || []).map((a: any) => ({ ...a, _type: 'article' })),
      users: (data.users || []).map((u: any) => ({ ...u, _type: 'user' })),
      prompts: (data.prompts || []).map((p: any) => ({ ...p, _type: 'prompt' }))
    }
    // Auto switch to tab with results
    const currentKey = activeTab.value as 'articles' | 'users' | 'prompts'
    if (allResults.value[currentKey].length === 0) {
      if (allResults.value.articles.length > 0) activeTab.value = 'articles'
      else if (allResults.value.users.length > 0) activeTab.value = 'users'
      else if (allResults.value.prompts.length > 0) activeTab.value = 'prompts'
    }
  } catch {
    allResults.value = { articles: [], users: [], prompts: [] }
  } finally {
    loading.value = false
  }
}

const discoverArticles = ref<any[]>([])

async function fetchDiscover() {
  try {
    const res = await get('/articles/hot', { limit: 10 })
    discoverArticles.value = res.data || []
  } catch { discoverArticles.value = [] }
}

onMounted(() => {
  if (route.query.q) {
    query.value = route.query.q as string
    doSearch()
  } else {
    fetchDiscover()
  }
})
</script>