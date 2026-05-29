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

    <div v-if="loading" class="text-center py-20 text-gray-400"><n-spin size="large" /></div>

    <template v-else-if="prompt">
      <div class="bg-white rounded-2xl border border-gray-100 p-6 lg:p-8 mb-6">
        <div class="flex items-center justify-between mb-4">
          <n-tag size="small" :bordered="false">{{ prompt.category }}</n-tag>
          <div class="flex items-center gap-2">
            <n-button quaternary size="small" @click="handleLike">
              {{ liked ? '❤️' : '🤍' }}
            </n-button>
            <n-button quaternary size="small" @click="handleFavorite">
              {{ favorited ? '⭐' : '☆' }}
            </n-button>
            <n-button quaternary size="small" @click="handleCopy">
              📋 复制
            </n-button>
          </div>
        </div>

        <h1 class="text-2xl font-bold text-dark mb-3">{{ prompt.title }}</h1>
        <p class="text-gray-500 mb-6">{{ prompt.description }}</p>

        <div class="flex items-center gap-4 mb-6 flex-wrap">
          <div class="flex items-center gap-2">
            <n-avatar :src="prompt.user?.avatar" round size="small" />
            <span class="text-sm text-gray-600">{{ prompt.user?.username }}</span>
          </div>
          <span class="text-xs text-gray-400">{{ formatDate(prompt.created_at) }}</span>
          <div class="flex items-center gap-1 text-sm">
            <span class="text-yellow-500">⭐</span>
            <span class="text-gray-600">{{ prompt.rating ? prompt.rating.toFixed(1) : '暂无评分' }}</span>
          </div>
          <span class="text-xs text-gray-400">使用 {{ formatCount(prompt.usage_count) }} 次</span>
          <span class="text-xs text-gray-400">❤️ {{ formatCount(prompt.like_count) }}</span>
        </div>

        <div class="bg-gray-50 rounded-xl p-6 mb-6">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-semibold text-dark">Prompt 内容</h3>
            <n-button size="tiny" quaternary @click="handleCopy">复制</n-button>
          </div>
          <pre class="text-sm text-gray-700 whitespace-pre-wrap font-mono">{{ prompt.content }}</pre>
        </div>

        <div v-if="prompt.tags?.length" class="flex items-center gap-2">
          <span class="text-xs text-gray-400">标签：</span>
          <n-tag v-for="tag in parseTags(prompt.tags)" :key="tag" size="tiny" :bordered="false">{{ tag }}</n-tag>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { get, post } from '@/api/request'
import { useUserStore } from '@/stores/user'
import type { Prompt } from '@/types/api'
import { formatDate, formatCount, parseTags } from '@/utils/helpers'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const userStore = useUserStore()

const prompt = ref<Prompt | null>(null)
const loading = ref(true)
const favorited = ref(false)
const liked = ref(false)

function handleCopy() {
  if (prompt.value?.content) {
    navigator.clipboard.writeText(prompt.value.content)
    message.success('已复制到剪贴板')
  }
}

async function handleFavorite() {
  try {
    await post(`/prompts/${route.params.id}/favorite`, {})
    favorited.value = !favorited.value
    message.success(favorited.value ? '已收藏' : '已取消收藏')
  } catch {
    message.error('操作失败')
  }
}

async function handleLike() {
  try {
    await post(`/prompts/${route.params.id}/like`, {})
    liked.value = !liked.value
    if (prompt.value) prompt.value.like_count += liked.value ? 1 : -1
    // 刷新用户信息（积分可能变化）
    userStore.fetchProfile()
  } catch {
    message.error('操作失败')
  }
}

async function fetchPrompt() {
  const id = Number(route.params.id)
  if (isNaN(id)) {
    message.error('Prompt 不存在')
    router.push('/prompt')
    return
  }
  loading.value = true
  try {
    const res = await get<Prompt>(`/prompts/${id}`)
    prompt.value = res.data
  } catch {
    message.error('获取 Prompt 失败')
    router.push('/prompt')
  } finally {
    loading.value = false
  }
}

onMounted(fetchPrompt)
</script>
