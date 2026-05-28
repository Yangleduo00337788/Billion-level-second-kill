<template>
  <div class="max-w-6xl mx-auto px-4 py-8">
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-bold text-dark">Prompt 社区</h1>
        <p class="text-sm text-gray-400 mt-1">发现和使用优质的 AI Prompt</p>
      </div>
      <router-link to="/prompt/create">
        <n-button type="primary" size="small">发布 Prompt</n-button>
      </router-link>
    </div>

    <div class="flex items-center gap-3 mb-6 overflow-x-auto pb-2">
      <n-button
        v-for="cat in categories"
        :key="cat"
        :type="selectedCategory === cat ? 'primary' : 'default'"
        size="tiny"
        :bordered="false"
        @click="selectedCategory = cat"
      >
        {{ cat }}
      </n-button>
    </div>

    <div class="relative mb-6">
      <n-input
        v-model:value="searchQuery"
        placeholder="搜索 Prompt..."
        size="large"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </template>
      </n-input>
    </div>

    <!-- Recommended Prompts -->
    <div v-if="recommendedPrompts.length > 0" class="mb-8">
      <h2 class="text-lg font-bold text-dark mb-4">精选推荐</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="item in recommendedPrompts"
          :key="item.recommend_id"
          class="bg-gradient-to-br from-purple-50 to-indigo-50 rounded-xl p-5 border border-purple-100 transition-all duration-200 hover:shadow-md hover:-translate-y-0.5 cursor-pointer"
          @click="router.push(`/prompt/${item.id}`)"
        >
          <div class="flex items-center gap-2 mb-3">
            <span class="px-2 py-0.5 text-xs bg-purple-100 text-purple-600 rounded-full">推荐</span>
            <span v-if="item.author" class="text-xs text-gray-400">{{ item.author }}</span>
          </div>
          <h3 class="font-semibold text-dark mb-2 line-clamp-2">{{ item.title }}</h3>
          <p v-if="item.summary" class="text-sm text-gray-500 line-clamp-2">{{ item.summary }}</p>
        </div>
      </div>
    </div>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="i in 6" :key="i" class="bg-white rounded-xl p-5 animate-pulse border border-gray-100">
        <div class="bg-gray-100 h-4 rounded w-2/3 mb-3"></div>
        <div class="bg-gray-100 h-3 rounded w-full mb-2"></div>
        <div class="bg-gray-100 h-3 rounded w-3/4"></div>
      </div>
    </div>

    <div v-else-if="prompts.length === 0" class="text-center py-20 text-gray-400">
      <p>暂无 Prompt</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="prompt in prompts"
        :key="prompt.id"
        class="bg-white rounded-xl p-5 border border-gray-100 transition-all duration-200 hover:shadow-md hover:-translate-y-0.5 cursor-pointer"
        @click="router.push(`/prompt/${prompt.id}`)"
      >
        <div class="flex items-center justify-between mb-3">
          <n-tag size="tiny" :bordered="false">{{ prompt.category }}</n-tag>
          <div class="flex items-center gap-1 text-xs text-gray-400">
            <span>⭐ {{ prompt.rating ? prompt.rating.toFixed(1) : '暂无评分' }}</span>
          </div>
        </div>
        <h3 class="font-semibold text-dark mb-2 line-clamp-1">{{ prompt.title }}</h3>
        <p class="text-sm text-gray-500 mb-4 line-clamp-2">{{ prompt.description }}</p>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <n-avatar :src="prompt.user?.avatar" round size="small" />
            <span class="text-xs text-gray-500">{{ prompt.user?.username }}</span>
          </div>
          <div class="flex items-center gap-3 text-xs text-gray-400">
            <span>使用 {{ formatCount(prompt.usage_count) }}</span>
            <span>❤️ {{ formatCount(prompt.like_count) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { get } from '@/api/request'
import { formatCount } from '@/utils/helpers'

const router = useRouter()

const selectedCategory = ref('全部')
const searchQuery = ref('')
const prompts = ref<any[]>([])
const loading = ref(false)
const recommendedPrompts = ref<Array<any>>([])

const categories = ref<string[]>(['全部'])

async function loadPromptCategories() {
  try {
    const res = await get('/articles/categories')
    const cats = res.data || []
    categories.value = ['全部', ...cats.map((c: any) => c.name)]
  } catch {}
}

loadPromptCategories()

async function fetchPrompts() {
  loading.value = true
  try {
    const params: any = { page: 1, page_size: 20 }
    if (selectedCategory.value !== '全部') params.category = selectedCategory.value
    if (searchQuery.value.trim()) params.q = searchQuery.value.trim()
    const res = await get('/prompts', params)
    prompts.value = res.data?.items || res.data || []
  } catch {
    prompts.value = []
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  fetchPrompts()
}

watch(selectedCategory, () => {
  fetchPrompts()
})

async function fetchRecommendedPrompts() {
  try {
    const res = await get<any>('/recommendations?position=prompt_recommend')
    recommendedPrompts.value = Array.isArray(res.data) ? res.data : []
  } catch {}
}

onMounted(() => {
  fetchPrompts()
  fetchRecommendedPrompts()
})
</script>
