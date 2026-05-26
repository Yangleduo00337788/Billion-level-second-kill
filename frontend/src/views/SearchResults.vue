<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <div class="mb-8">
      <n-input
        v-model:value="query"
        placeholder="搜索文章、用户、Prompt..."
        size="large"
        clearable
        @keyup.enter="doSearch"
      >
        <template #prefix>
          <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </template>
      </n-input>
    </div>

    <div class="flex items-center gap-1 mb-6 border-b border-gray-100">
      <n-button
        v-for="tab in tabs"
        :key="tab.key"
        :type="activeTab === tab.key ? 'primary' : 'default'"
        :bordered="false"
        size="small"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </n-button>
    </div>

    <div v-if="!query && !hasSearched" class="text-center py-20 text-gray-400">
      <p>输入关键词开始搜索</p>
    </div>

    <div v-else-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="bg-white rounded-xl p-5 animate-pulse border border-gray-100">
        <div class="bg-gray-100 h-4 rounded w-2/3 mb-3"></div>
        <div class="bg-gray-100 h-3 rounded w-full mb-2"></div>
        <div class="bg-gray-100 h-3 rounded w-3/4"></div>
      </div>
    </div>

    <div v-else-if="results.length === 0 && hasSearched" class="text-center py-20 text-gray-400">
      <p>没有找到相关结果</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="item in results"
        :key="item.id"
        class="bg-white rounded-xl border border-gray-100 p-5 cursor-pointer hover:shadow-sm transition-shadow"
        @click="goToDetail(item)"
      >
        <h3 class="font-semibold text-dark mb-2" v-html="highlightText(item.title)"></h3>
        <p class="text-sm text-gray-500 line-clamp-2 mb-3" v-html="highlightText(item.summary || item.description || '')"></p>
        <div class="flex items-center gap-3 text-xs text-gray-400">
          <span>{{ item.user?.username }}</span>
          <span>{{ formatDate(item.created_at) }}</span>
          <span v-if="item.view_count !== undefined">{{ formatCount(item.view_count) }} 阅读</span>
          <span v-if="item.usage_count !== undefined">使用 {{ formatCount(item.usage_count) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get } from '@/api/request'
import { formatDate, formatCount } from '@/utils/helpers'

const route = useRoute()
const router = useRouter()

const query = ref('')
const activeTab = ref('articles')
const results = ref<any[]>([])
const loading = ref(false)
const hasSearched = ref(false)

const tabs = [
  { key: 'articles', label: '文章' },
  { key: 'users', label: '用户' },
  { key: 'prompts', label: 'Prompt' }
]

function escapeHtml(text: string): string {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

function highlightText(text: string): string {
  if (!text || !query.value) return text || ''
  const safe = escapeHtml(text)
  const q = query.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return safe.replace(new RegExp(`(${q})`, 'gi'), '<mark class="bg-yellow-100 text-yellow-900 px-0.5 rounded">$1</mark>')
}

function goToDetail(item: any) {
  const type = item._type || activeTab.value
  if (type === 'article') router.push(`/article/${item.id}`)
  else if (type === 'prompt') router.push(`/prompt/${item.id}`)
  else if (type === 'user') router.push(`/user/${item.id}`)
}

async function doSearch() {
  if (!query.value.trim()) return
  loading.value = true
  hasSearched.value = true
  try {
    const res = await get('/search', { q: query.value, page: 1, size: 20 })
    const data = res.data as any
    const allResults: any[] = []

    if (activeTab.value === 'articles') {
      if (data.articles) allResults.push(...data.articles.map((a: any) => ({ ...a, _type: 'article' })))
    }
    if (activeTab.value === 'users') {
      if (data.users) allResults.push(...data.users.map((u: any) => ({ ...u, _type: 'user' })))
    }
    if (activeTab.value === 'prompts') {
      if (data.prompts) allResults.push(...data.prompts.map((p: any) => ({ ...p, _type: 'prompt' })))
    }

    results.value = allResults
  } catch {
    results.value = []
  } finally {
    loading.value = false
  }
}

watch(activeTab, () => {
  if (hasSearched.value) doSearch()
})

onMounted(() => {
  if (route.query.q) {
    query.value = route.query.q as string
    doSearch()
  }
})
</script>
